package main

import (
   "bytes"
   "crypto/tls"
   "encoding/json"
   "fmt"
   "github.com/89z/parse/crypto"
   "io"
   "io/ioutil"
   "net"
   "net/http"
   "net/http/httptrace"
   "net/url"
   "strings"
   "sync/atomic"
   "time"
   stdcontext "context"
)

type spyConn struct {
   net.Conn
}

// "\x16\x03\x01\x02\x00\x01\x00\x01\xfc"
// 769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0
func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   if p[0] == 0x16 {
      fmt.Println("Handshake")
      for _, hand := range crypto.Handshakes(p) {
         hello, err := crypto.ParseHandshake(hand)
         if err == nil {
            ja3, err := hello.FormatJA3()
            if err == nil {
               fmt.Println(ja3)
            }
         }
      }
   }
   return n, err
}

const defaultTargetConnectTimeout   = 5 * time.Second

// Canned HTTP responses
var (
   badGateway = fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
   internalErr = "PROXY_CHANNEL_INTERNAL_ERR"
   tunnelEstablishedResponseLine = []byte(fmt.Sprintf("HTTP/1.1 %d Connection established\r\n\r\n", http.StatusOK))
)

// ProxyError specifies all the possible errors that can occur due to this proxy's behavior,
// which does not include the behavior of parent proxies.
type proxyError struct {
	ErrType string `json:"errType"`
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

const NormalMode = iota

////////////////////////////////////////////////////////////////////////////////

// Proxy is a struct that implements ServeHTTP() method
type proxy struct {
   clientConnNum int32
   transport     *http.Transport
   mode          int
}

func NewProxy(hconf *http.Transport) *proxy {
   p := &proxy{}
   if hconf == nil {
      p.transport = &http.Transport{
         TLSClientConfig: &tls.Config{
            // No need to verify because as a proxy we don't care
            InsecureSkipVerify: true,
         },
         DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
            DualStack: true,
         }).DialContext,
         MaxIdleConns:          100,
         MaxIdleConnsPerHost:   10,
         IdleConnTimeout:       90 * time.Second,
         TLSHandshakeTimeout:   10 * time.Second,
         ExpectContinueTimeout: 1 * time.Second,
         ProxyConnectHeader:    make(http.Header),
      }
   } else {
      p.transport = hconf
      p.transport.ProxyConnectHeader = make(http.Header)
   }
   return p
}

// ServeHTTP .
func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
   fmt.Println(req.Header)
   if req.URL.Host == "" {
      req.URL.Host = req.Host
   }
   atomic.AddInt32(&p.clientConnNum, 1)
   defer func() {
      atomic.AddInt32(&p.clientConnNum, -1)
   }()
   ctx := &context{
      data:       make(map[interface{}]interface{}),
      req:        req,
   }
   if ctx.abort {
      ctx.setContextErrType(ConnectFail)
      return
   }
   if ctx.abort {
      ctx.setContextErrType(AuthFail)
      return
   }
   // NormalMode: This proxy will forward requests to parent proxy, and return
   // whatever it gets from parent proxy back to requestor.
   switch p.mode {
   case NormalMode:
      if ctx.req.Method == http.MethodConnect {
         h := ctx.req.Header.Get("MITM")
         if h == "Enabled" {
            ctx.mitm = true
         } else {
            p.proxyTunnel(ctx, rw)
         }
      } else {
         p.proxyHTTP(ctx, rw)
      }
   }
}

// ClientConnNum gets the Client
func (p *proxy) ClientConnNum() int32 {
	return atomic.LoadInt32(&p.clientConnNum)
}

func writeProxyErrorToResponseBody(ctx *context, respWriter io.Writer, httpcode int32, msg string, optionalPrefix string) {
	if optionalPrefix != "" {
		m, _ := respWriter.Write([]byte(optionalPrefix))
		ctx.respLength += int64(m)
	}
	pe := &proxyError{
		ErrType: internalErr,
		ErrCode: httpcode,
		ErrMsg:  msg,
	}
	errJSON, err := json.Marshal(pe)
	if err != nil {
		panic(fmt.Errorf("jason marshal failed"))
	}
	n, _ := respWriter.Write(errJSON)
	ctx.respLength += int64(n)
}

func (p *proxy) proxyHTTP(ctx *context, rw http.ResponseWriter) {
   fmt.Println(ctx.req.Header)
   ctx.req.URL.Scheme = "http"
   p.DoRequest(ctx, rw, func(resp *http.Response, err error) {
      if err != nil {
         fmt.Printf("proxyHTTP %s forward request failed: %s", ctx.req.URL, err)
         rw.WriteHeader(http.StatusBadGateway)
         writeProxyErrorToResponseBody(ctx, rw, http.StatusBadGateway, fmt.Sprintf("proxyHTTP %s forward request failed: %s", ctx.req.URL, err), "")
         ctx.setContextErrorWithType(err, HTTPDoRequestFail)
         return
      }
      defer resp.Body.Close()
      CopyHeader(rw.Header(), resp.Header)
      rw.WriteHeader(resp.StatusCode)
      written, err := io.Copy(rw, resp.Body)
      ctx.respLength += written
      if err != nil {
         fmt.Printf("proxyHTTP %s write client failed: %s", ctx.req.URL, err)
         ctx.setContextErrorWithType(err, HTTPWriteClientFail)
         return
      }
   })
}

func (p *proxy) proxyTunnel(ctx *context, rw http.ResponseWriter) {
   if ctx.abort {
      ctx.setContextErrType(ParentProxyFail)
      return
   }
   clientConn, err := hijacker(rw)
   if err != nil {
      fmt.Printf("proxyTunnel hijack client connection failed: %s", err)
      rw.WriteHeader(http.StatusBadGateway)
      writeProxyErrorToResponseBody(ctx, rw, http.StatusBadGateway, fmt.Sprintf("proxyTunnel hijack client connection failed: %s", err), "")
      ctx.setContextErrorWithType(err, TunnelHijackClientConnFail)
      return
   }
   ctx.hijack = true
   defer func() {
      err := clientConn.Close()
      if err != nil {
         fmt.Printf("defer client close err: %s", err)
      } else {
         fmt.Println("defer client close done")
      }
   }()
   fmt.Println(ctx.req.Header)
   targetAddr := ctx.req.URL.Host
   targetConn, err := net.DialTimeout("tcp", targetAddr, defaultTargetConnectTimeout)
   if ctx.abort {
      ctx.setContextErrType(BeforeResponseFail)
      return
   }
   if err != nil {
      fmt.Printf("proxyTunnel %s dial remote server failed: %s", ctx.req.URL.Host, err)
      writeProxyErrorToResponseBody(ctx, clientConn, http.StatusBadGateway, fmt.Sprintf("proxyTunnel %s dial remote server failed: %s", ctx.req.URL.Host, err), badGateway)
      ctx.setContextErrorWithType(err, TunnelDialRemoteServerFail)
      return
   }
   defer func() {
      err := targetConn.Close()
      if err != nil {
         fmt.Printf("defer target close err: %s", err)
      } else {
         fmt.Println("defer target close done")
      }
   }()
   _, err = clientConn.Write(tunnelEstablishedResponseLine)
   if err != nil {
      fmt.Printf("proxyTunnel %s write message failed: %s", ctx.req.URL.Host, err)
      ctx.setContextErrorWithType(err, TunnelWriteEstRespFail)
      return
   }
   transfer(ctx, clientConn, targetConn)
}

// transfer does two-way forwarding through connections
func transfer(ctx *context, clientConn net.Conn, targetConn net.Conn, parentProxy ...string) {
   go func() {
      written1, err1 := io.Copy(clientConn, targetConn)
      if err1 != nil {
         fmt.Printf("io.Copy write clientConn failed: %s", err1)
         if len(parentProxy) <= 1 {
            if len(parentProxy) == 0 {
               ctx.setContextErrorWithType(err1, TunnelWriteClientConnFinish)
            } else {
               ctx.setPoolContextErrorWithType(err1, TunnelWriteClientConnFinish, parentProxy[0])
            }
         }
      }
      ctx.respLength += written1
      clientConn.Close()
      targetConn.Close()
   }()
   spy := spyConn{clientConn}
   written2, err2 := io.Copy(targetConn, spy)
   if err2 != nil {
   fmt.Printf("io.Copy write targetConn failed: %s", err2)
   if len(parentProxy) <= 1 {
   if len(parentProxy) == 0 {
   ctx.setContextErrorWithType(err2, TunnelWriteTargetConnFinish)
   } else {
   ctx.setPoolContextErrorWithType(err2, TunnelWriteTargetConnFinish, parentProxy[0])
   }
   }
   }
   ctx.reqLength += written2
   targetConn.Close()
   clientConn.Close()
}

// DoRequest makes a request to remote server as a clent through given proxy,
// and calls responseFunc before returning the response.
// The "conn" is needed when it comes to https request, and only one conn is accepted.
func (p *proxy) DoRequest(ctx *context, rw http.ResponseWriter, responseFunc func(*http.Response, error), conn ...interface{}) {
   if len(conn) > 1 {
      return
   }
   if ctx.data == nil {
      ctx.data = make(map[interface{}]interface{})
   }
   if ctx.abort {
      ctx.setContextErrType(BeforeRequestFail)
      return
   }
   newReq := new(http.Request)
   *newReq = *ctx.req
   fmt.Println(newReq.Header)
   newReq.Header = CloneHeader(newReq.Header)
   // When server reads http request it sets req.Close to true if "Connection"
   // header contains "close".
   // https://github.com/golang/go/blob/master/src/net/http/request.go#L1080
   // Later, transfer.go adds "Connection: close" back when req.Close is true
   // https://github.com/golang/go/blob/master/src/net/http/transfer.go#L275
   // That's why tests that checks "Connection: close" removal fail
   if newReq.Header.Get("Connection") == "close" {
      newReq.Close = false
   }
   removeMITMHeaders(newReq.Header)
   removeConnectionHeaders(newReq.Header)
   removeHopHeaders(newReq.Header)
   var parentProxyURL *url.URL
   var err error
   if ctx.abort {
      ctx.setContextErrType(ParentProxyFail)
      return
   }
   type CtxKey int
   var pkey CtxKey = 0
   fakeCtx := stdcontext.WithValue(newReq.Context(), pkey, parentProxyURL)
   newReq = newReq.Clone(fakeCtx)
   ctx.reqLength += newReq.ContentLength
   tr := p.transport
   tr.Proxy = func(req *http.Request) (*url.URL, error) {
      ctx := req.Context()
      pURL := ctx.Value(pkey).(*url.URL)
      trace := &httptrace.ClientTrace{
      GotConn: func(connInfo httptrace.GotConnInfo) {
      fmt.Printf("Got conn: %+v", connInfo)
      },
      DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
      fmt.Printf("DNS done, info: %+v", dnsInfo)
      },
      GotFirstResponseByte: func() {
      fmt.Printf("GotFirstResponseByte: %+v", time.Now())
      },
      }
      req.Clone(httptrace.WithClientTrace(stdcontext.Background(), trace))
      return pURL, err
   }
   resp, err := tr.RoundTrip(newReq)
   if ctx.abort {
      ctx.setContextErrType(BeforeResponseFail)
      return
   }
   if err == nil {
      removeConnectionHeaders(resp.Header)
      removeHopHeaders(resp.Header)
   }
   responseFunc(resp, err)
}

// hijacker gets the underlying connection of an http.ResponseWriter
func hijacker(rw http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		return nil, fmt.Errorf("hijacker is not supported")
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("hijacker failed: %s", err)
	}

	return conn, nil
}

var hopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func removeConnectionHeaders(h http.Header) {
	if c := h.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				h.Del(f)
			}
		}
	}
}

func removeHopHeaders(h http.Header) {
	for _, item := range hopHeaders {
		if h.Get(item) != "" {
			h.Del(item)
		}
	}
}

func removeMITMHeaders(h http.Header) {
	if c := h.Get("MITM"); c != "" {
		h.Del("MITM")
	}
}

// CopyHeader shallow copy.
func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// CloneHeader deep copy.
func CloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// CloneBody deep copy.
func CloneBody(b io.ReadCloser) (r io.ReadCloser, body []byte, err error) {
	if b == nil {
		return http.NoBody, nil, nil
	}
	body, err = ioutil.ReadAll(b)
	if err != nil {
		return http.NoBody, nil, err
	}
	r = ioutil.NopCloser(bytes.NewReader(body))

	return r, body, nil
}