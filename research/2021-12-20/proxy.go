package proxychannel

import (
   "bytes"
   "context"
   "crypto/tls"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "github.com/jmcvetta/randutil"
   "io"
   "io/ioutil"
   "net"
   "net/http"
   "net/http/httptrace"
   "net/url"
   "strings"
   "sync/atomic"
   "time"
   "github.com/89z/parse/crypto"
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

// Default timeout values
const (
	defaultTargetConnectTimeout   = 5 * time.Second
	defaultTargetReadWriteTimeout = 30 * time.Second
	defaultClientReadWriteTimeout = 30 * time.Second
)

const defaultHTTPResponsePeekSize int = 4096

// Canned HTTP responses
var tunnelEstablishedResponseLine = []byte(fmt.Sprintf("HTTP/1.1 %d Connection established\r\n\r\n", http.StatusOK))
var badGateway = fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
var tooManyRequests = fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))

var internalErr = "PROXY_CHANNEL_INTERNAL_ERR"

// ProxyError specifies all the possible errors that can occur due to this proxy's behavior,
// which does not include the behavior of parent proxies.
type ProxyError struct {
	ErrType string `json:"errType"`
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

// TunnelConn .
type TunnelConn struct {
	Client net.Conn
	Target net.Conn
}

// TunnelInfo .
type TunnelInfo struct {
	Client      net.Conn
	Target      net.Conn
	Err         error
	ParentProxy *url.URL
	Pool        ConnPool
}

// ResponseInfo .
type ResponseInfo struct {
	Resp        *http.Response
	Err         error
	ParentProxy *url.URL
	Pool        ConnPool
}

// ResponseWrapper is simply a wrapper for http.Response and error.
type ResponseWrapper struct {
	Resp *http.Response
	Err  error
}

// ConnWrapper .
type ConnWrapper struct {
	Conn net.Conn
	Err  error
}

// Below are the modes supported.
const (
	NormalMode = iota
	ConnPoolMode
)

func makeTunnelRequestLine(addr string) string {
	return fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", addr)
}

func makeTunnelRequestWithAuth(ctx *Context, parentProxyURL *url.URL, targetConn net.Conn) error {
	connectReq := &http.Request{
		Proto:      ctx.Req.Proto,
		ProtoMajor: ctx.Req.ProtoMajor,
		ProtoMinor: ctx.Req.ProtoMinor,
		Method:     "CONNECT",
		URL:        &url.URL{Opaque: ctx.Req.URL.Host},
		Host:       ctx.Req.URL.Host,
		Header:     CloneHeader(ctx.Req.Header),
	}
	if connectReq.Proto == "HTTP/1.0" {
		connectReq.Header.Del("Connection")
	}
	u := parentProxyURL.User
	if u != nil {
		username := u.Username()
		password, _ := u.Password()
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		connectReq.Header.Set("Proxy-Authorization", basicAuth)
	}
	return connectReq.Write(targetConn)
}

// Proxy is a struct that implements ServeHTTP() method
type Proxy struct {
	delegate      Delegate
	clientConnNum int32
	decryptHTTPS  bool
	//cert          *cert.Certificate
	transport     *http.Transport
	mode          int
}

var _ http.Handler = &Proxy{}

func NewProxy(hconf *HandlerConfig, em *ExtensionManager) *Proxy {
	p := &Proxy{}

	if hconf.Delegate == nil {
		p.delegate = &DefaultDelegate{}
	} else {
		p.delegate = hconf.Delegate
	}
	p.delegate.SetExtensionManager(em)

	//p.cert = cert.NewCertificate(hconf.CertCache)

	if hconf.Transport == nil {
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
		p.transport = hconf.Transport
		p.transport.ProxyConnectHeader = make(http.Header)
	}
	p.transport.DisableKeepAlives = hconf.DisableKeepAlive
	p.mode = hconf.Mode
	if p.mode == ConnPoolMode {
		p.transport.ProxyConnectHeader.Set("MITM", "Enabled")
	}
	return p
}

// ServeHTTP .
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
         fmt.Println(req.Header)
         if req.URL.Host == "" {
		req.URL.Host = req.Host
	}
	atomic.AddInt32(&p.clientConnNum, 1)
	defer func() {
		atomic.AddInt32(&p.clientConnNum, -1)
	}()
	ctx := &Context{
		Req:        req,
		Data:       make(map[interface{}]interface{}),
		Hijack:     false,
		MITM:       false,
		ReqLength:  0,
		RespLength: 0,
		Closed:     false,
	}
	defer p.delegate.Finish(ctx, rw)
	p.delegate.Connect(ctx, rw)
	if ctx.abort {
		ctx.SetContextErrType(ConnectFail)
		return
	}
	p.delegate.Auth(ctx, rw)
	if ctx.abort {
		ctx.SetContextErrType(AuthFail)
		return
	}

	// NormalMode:
	// This proxy will forward requests to parent proxy, and return whatever it gets
	// from parent proxy back to requestor.

	// ConnPoolMode:
	// This proxy chooses a TCP connection by given probability from the ConnPool,
	// which is specified by p.delegate.GetConnPool(ctx).
	// If this proxy fails to connect parent proxy or gets a response body in JSON
	// format that has an "ErrType" of "PROXY_CHANNEL_INTERNAL_ERR"(especially when
	// the parent proxy is also a proxychannel instance), it retries the proxy request
	// with another proxy chosen from ConnPool by given probability.
	// The retry goes on until any parent proxy returns a 200 response code or every
	// connection has been chosen.
	switch p.mode {
	case NormalMode:
               if ctx.Req.Method == http.MethodConnect {
			h := ctx.Req.Header.Get("MITM")
			if h == "Enabled" {
				ctx.MITM = true
			} else {
				p.proxyTunnel(ctx, rw)
			}
		} else {
			   p.proxyHTTP(ctx, rw)
		}
	}
}

// ClientConnNum gets the Client
func (p *Proxy) ClientConnNum() int32 {
	return atomic.LoadInt32(&p.clientConnNum)
}

// WriteProxyErrorToResponseBody is the standard function to call when errors occur due to this proxy's behavior,
// which does not include the behavior of parent proxies.
func WriteProxyErrorToResponseBody(ctx *Context, respWriter Writer, httpcode int32, msg string, optionalPrefix string) {
	if optionalPrefix != "" {
		m, _ := respWriter.Write([]byte(optionalPrefix))
		ctx.RespLength += int64(m)
	}
	pe := &ProxyError{
		ErrType: internalErr,
		ErrCode: httpcode,
		ErrMsg:  msg,
	}
	errJSON, err := json.Marshal(pe)
	if err != nil {
		panic(fmt.Errorf("jason marshal failed"))
	}
	n, _ := respWriter.Write(errJSON)
	ctx.RespLength += int64(n)
}

func (p *Proxy) proxyHTTP(ctx *Context, rw http.ResponseWriter) {
      fmt.Println(ctx.Req.Header)
      ctx.Req.URL.Scheme = "http"
	p.DoRequest(ctx, rw, func(resp *http.Response, err error) {
		if err != nil {
			Logger.Errorf("proxyHTTP %s forward request failed: %s", ctx.Req.URL, err)
			rw.WriteHeader(http.StatusBadGateway)
			WriteProxyErrorToResponseBody(ctx, rw, http.StatusBadGateway, fmt.Sprintf("proxyHTTP %s forward request failed: %s", ctx.Req.URL, err), "")
			ctx.SetContextErrorWithType(err, HTTPDoRequestFail)
			return
		}

		defer resp.Body.Close()
		p.delegate.DuringResponse(ctx, resp)

		CopyHeader(rw.Header(), resp.Header)
		rw.WriteHeader(resp.StatusCode)

		written, err := io.Copy(rw, resp.Body)
		ctx.RespLength += written
		if err != nil {
			Logger.Errorf("proxyHTTP %s write client failed: %s", ctx.Req.URL, err)
			ctx.SetContextErrorWithType(err, HTTPWriteClientFail)
			return
		}
	})
}

func (p *Proxy) proxyTunnel(ctx *Context, rw http.ResponseWriter) {
	parentProxyURL, err := p.delegate.ParentProxy(ctx, rw)
	if ctx.abort {
		ctx.SetContextErrType(ParentProxyFail)
		return
	}
	clientConn, err := hijacker(rw)
	if err != nil {
		Logger.Errorf("proxyTunnel hijack client connection failed: %s", err)
		rw.WriteHeader(http.StatusBadGateway)
		WriteProxyErrorToResponseBody(ctx, rw, http.StatusBadGateway, fmt.Sprintf("proxyTunnel hijack client connection failed: %s", err), "")
		ctx.SetContextErrorWithType(err, TunnelHijackClientConnFail)
		return
	}
	ctx.Hijack = true
	defer func() {
		err := clientConn.Close()
		if err != nil {
			Logger.Infof("defer client close err: %s", err)
		} else {
			Logger.Infof("defer client close done")
		}
	}()
	// defer clientConn.Close()
      fmt.Println(ctx.Req.Header)
	targetAddr := ctx.Req.URL.Host
	if parentProxyURL != nil {
		targetAddr = parentProxyURL.Host
	}

	targetConn, err := net.DialTimeout("tcp", targetAddr, defaultTargetConnectTimeout)

	connWrapper := &ConnWrapper{
		Conn: targetConn,
		Err:  err,
	}
	p.delegate.BeforeResponse(ctx, connWrapper)
	if ctx.abort {
		ctx.SetContextErrType(BeforeResponseFail)
		return
	}
	if err != nil {
		Logger.Errorf("proxyTunnel %s dial remote server failed: %s", ctx.Req.URL.Host, err)
		WriteProxyErrorToResponseBody(ctx, clientConn, http.StatusBadGateway, fmt.Sprintf("proxyTunnel %s dial remote server failed: %s", ctx.Req.URL.Host, err), badGateway)
		ctx.SetContextErrorWithType(err, TunnelDialRemoteServerFail)
		return
	}
	// defer targetConn.Close()
	defer func() {
		err := targetConn.Close()
		if err != nil {
			Logger.Infof("defer target close err: %s", err)
		} else {
			Logger.Infof("defer target close done")
		}
	}()
	p.delegate.DuringResponse(ctx, &TunnelConn{Client: clientConn, Target: targetConn}) // targetConn could be closed in this method
	if parentProxyURL == nil {
		_, err = clientConn.Write(tunnelEstablishedResponseLine)
		if err != nil {
			Logger.Errorf("proxyTunnel %s write message failed: %s", ctx.Req.URL.Host, err)
			ctx.SetContextErrorWithType(err, TunnelWriteEstRespFail)
			return
		}
	} else {
		err := makeTunnelRequestWithAuth(ctx, parentProxyURL, targetConn)
		if err != nil {
			Logger.Errorf("proxyTunnel %s make connect request to remote failed: %s", ctx.Req.URL.Host, err)
			WriteProxyErrorToResponseBody(ctx, clientConn, http.StatusBadGateway, fmt.Sprintf("proxyTunnel %s make connect request to remote failed: %s", ctx.Req.URL.Host, err), badGateway)
			ctx.SetContextErrorWithType(err, TunnelConnectRemoteFail)
			return
		}
	}
	transfer(ctx, clientConn, targetConn)
}

// transfer does two-way forwarding through connections
func transfer(ctx *Context, clientConn net.Conn, targetConn net.Conn, parentProxy ...string) {
   go func() {
      written1, err1 := io.Copy(clientConn, targetConn)
      if err1 != nil {
         Logger.Errorf("io.Copy write clientConn failed: %s", err1)
         if len(parentProxy) <= 1 {
            if len(parentProxy) == 0 {
               ctx.SetContextErrorWithType(err1, TunnelWriteClientConnFinish)
            } else {
               ctx.SetPoolContextErrorWithType(err1, TunnelWriteClientConnFinish, parentProxy[0])
            }
         }
      }
      ctx.RespLength += written1
      clientConn.Close()
      targetConn.Close()
   }()
   
   //written2, err2 := io.Copy(targetConn, clientConn)
   
   spy := spyConn{clientConn}
   written2, err2 := io.Copy(targetConn, spy)
   
	if err2 != nil {
		Logger.Errorf("io.Copy write targetConn failed: %s", err2)
		if len(parentProxy) <= 1 {
			if len(parentProxy) == 0 {
				ctx.SetContextErrorWithType(err2, TunnelWriteTargetConnFinish)
			} else {
				ctx.SetPoolContextErrorWithType(err2, TunnelWriteTargetConnFinish, parentProxy[0])
			}
		}
	}
	ctx.ReqLength += written2
	targetConn.Close()
	clientConn.Close()
}


// DoRequest makes a request to remote server as a clent through given proxy,
// and calls responseFunc before returning the response.
// The "conn" is needed when it comes to https request, and only one conn is accepted.
func (p *Proxy) DoRequest(ctx *Context, rw http.ResponseWriter, responseFunc func(*http.Response, error), conn ...interface{}) {
	if len(conn) > 1 {
		return
	}
	var clientConn *tls.Conn
	if len(conn) == 1 {
		c := conn[0]
		clientConn, _ = c.(*tls.Conn)
	}

	if ctx.Data == nil {
		ctx.Data = make(map[interface{}]interface{})
	}
	p.delegate.BeforeRequest(ctx)
	if ctx.abort {
		ctx.SetContextErrType(BeforeRequestFail)
		return
	}
	newReq := new(http.Request)
	*newReq = *ctx.Req
      fmt.Println(newReq.Header)
	newReq.Header = CloneHeader(newReq.Header)
	// When server reads http request it sets req.Close to true if
	// "Connection" header contains "close".
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

	// p.transport.ForceAttemptHTTP2 = true // for HTTP/2 test
	var parentProxyURL *url.URL
	var err error
	if ctx.Hijack {
		parentProxyURL, err = p.delegate.ParentProxy(ctx, clientConn)
	} else {
		parentProxyURL, err = p.delegate.ParentProxy(ctx, rw)
	}
	if ctx.abort {
		ctx.SetContextErrType(ParentProxyFail)
		return
	}

	type CtxKey int
	var pkey CtxKey = 0
	fakeCtx := context.WithValue(newReq.Context(), pkey, parentProxyURL)
	newReq = newReq.Clone(fakeCtx)

	ctx.ReqLength += newReq.ContentLength
	// dump, dumperr := httputil.DumpRequestOut(newReq, true)
	// if dumperr != nil {
	// 	Logger.Errorf("DumpRequestOut failed %s", dumperr)
	// } else {
	// 	ctx.ReqLength += int64(len(dump))
	// }

	tr := p.transport
	tr.Proxy = func(req *http.Request) (*url.URL, error) {
		ctx := req.Context()
		pURL := ctx.Value(pkey).(*url.URL)
		// req = req.Clone(context.Background())
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				Logger.Infof("Got conn: %+v", connInfo)
			},
			DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
				Logger.Infof("DNS done, info: %+v", dnsInfo)
			},
			GotFirstResponseByte: func() {
				Logger.Infof("GotFirstResponseByte: %+v", time.Now())
			},
		}
		req = req.Clone(httptrace.WithClientTrace(context.Background(), trace))
		return pURL, err
	}

	resp, err := tr.RoundTrip(newReq)

	respWrapper := &ResponseWrapper{
		Resp: resp,
		Err:  err,
	}

	p.delegate.BeforeResponse(ctx, respWrapper)
	if ctx.abort {
		ctx.SetContextErrType(BeforeResponseFail)
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

func (p *Proxy) proxyHTTPWithConnPool(ctx *Context, rw http.ResponseWriter) {
	ctx.Req.URL.Scheme = "http"
	if ctx.Data == nil {
		ctx.Data = make(map[interface{}]interface{})
	}
	newReq := new(http.Request)
	*newReq = *ctx.Req
	newReq.Header = CloneHeader(newReq.Header)
	// When server reads http request it sets req.Close to true if
	// "Connection" header contains "close".
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

	poolChoices, err := p.delegate.GetConnPool(ctx)
	if err != nil {
		Logger.Errorf("proxyHTTPWithConnPool %s GetConnPool failed: %s", ctx.Req.URL.Host, err)
		rw.WriteHeader(http.StatusBadGateway)
		WriteProxyErrorToResponseBody(ctx, rw, http.StatusBadGateway, fmt.Sprintf("proxyHTTPWithConnPool %s GetConnPool failed: %s", ctx.Req.URL.Host, err), "")
		ctx.SetPoolContextErrorWithType(err, PoolGetParentProxyFail)
		return
	}

	work := false
	for range poolChoices {
		// pool is not actully used to connect parent proxy,
		// it's just used to show whether a connection to it is performing well.
		// I know it's weird, and it will be fixed in the future.

		choice, err := randutil.WeightedChoice(poolChoices)
		pool := choice.Item.(ConnPool)
		parentProxyURL := pool.GetRemoteAddrURL()
		proxyTag := pool.GetTag()
		for i := range poolChoices {
			pl := poolChoices[i].Item.(ConnPool)
			if pl.GetTag() == proxyTag {
				poolChoices[i].Weight = 0
				break
			}
		}

		type CtxKey int
		var pkey CtxKey = 0
		fakeCtx := context.WithValue(newReq.Context(), pkey, parentProxyURL)
		newReq = newReq.Clone(fakeCtx)

		ctx.ReqLength += newReq.ContentLength
		// dump, dumperr := httputil.DumpRequestOut(newReq, true)
		// if dumperr != nil {
		// 	Logger.Errorf("DumpRequestOut failed %s", dumperr)
		// } else {
		// 	ctx.ReqLength += int64(len(dump))
		// }

		tr := p.transport
		tr.Proxy = func(req *http.Request) (*url.URL, error) {
			ctx := req.Context()
			pURL := ctx.Value(pkey).(*url.URL)
			// req = req.Clone(context.Background())
			trace := &httptrace.ClientTrace{
				GotConn: func(connInfo httptrace.GotConnInfo) {
					Logger.Infof("Got conn: %+v", connInfo)
				},
				DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
					Logger.Infof("DNS done, info: %+v", dnsInfo)
				},
				GotFirstResponseByte: func() {
					Logger.Infof("GotFirstResponseByte: %+v", time.Now())
				},
			}
			req = req.Clone(httptrace.WithClientTrace(context.Background(), trace))
			return pURL, err
		}

		resp, err := tr.RoundTrip(newReq)
		p.delegate.BeforeResponse(ctx, &ResponseInfo{
			Resp:        resp,
			Err:         err,
			ParentProxy: parentProxyURL,
			Pool:        pool,
		})
		if ctx.abort {
			ctx.SetPoolContextErrorWithType(nil, BeforeRequestFail)
			return
		}

		if err != nil {
			Logger.Errorf("proxyHTTPWithConnPool %s RoundTrip failed: %s", ctx.Req.URL, err)
			ctx.SetPoolContextErrorWithType(err, PoolRoundTripFail, proxyTag)
			continue
		}
		removeConnectionHeaders(resp.Header)
		removeHopHeaders(resp.Header)

		// defer resp.Body.Close() is not used as it's in a loop.
		p.delegate.DuringResponse(ctx, resp)

		buf := make([]byte, defaultHTTPResponsePeekSize+1) // max acceptable size is defaultHTTPResponsePeekSize bytes.
		n, err := io.ReadFull(resp.Body, buf)
		switch err {
		case nil:
		case io.ErrUnexpectedEOF:
		case io.EOF:
			n = 0
		// Only errors that are not listed above will be treated as real error
		default:
			Logger.Errorf("proxyHTTPWithConnPool %s ReadFull failed: %s", ctx.Req.URL, err)
			ctx.SetPoolContextErrorWithType(err, PoolReadRemoteFail, proxyTag)
			resp.Body.Close()
			continue
		}
		buf = buf[:n]
		if resp.StatusCode != http.StatusTooManyRequests {
			if resp.StatusCode == http.StatusOK || !strings.Contains(string(buf), internalErr) {
				// No need to retry, just return what we get to rw.
				work = true
				CopyHeader(rw.Header(), resp.Header)
				rw.WriteHeader(resp.StatusCode)
				m, err := rw.Write(buf)
				ctx.RespLength += int64(m)
				if err != nil || n != m {
					if err != nil {
						Logger.Errorf("proxyHTTPWithConnPool %s first part write client failed:", ctx.Req.URL, err)
						ctx.SetPoolContextErrorWithType(err, PoolWriteClientFail, proxyTag)
					} else {
						Logger.Errorf("proxyHTTPWithConnPool %s partial write, read: %d, write: %d", ctx.Req.URL, n, m)
						ctx.SetPoolContextErrorWithType(fmt.Errorf("proxyHTTPWithConnPool %s partial write, read: %d, write: %d", ctx.Req.URL, n, m), PoolWriteClientFail, proxyTag)
					}
					resp.Body.Close()
					break
				}
				written, err := io.Copy(rw, resp.Body)
				ctx.RespLength += written
				if err != nil {
					Logger.Errorf("proxyHTTPWithConnPool %s write client failed: %s", ctx.Req.URL, err)
					ctx.SetPoolContextErrorWithType(err, PoolWriteClientFail, proxyTag)
					resp.Body.Close()
					break
				}
				ctx.SetPoolContextErrorWithType(fmt.Errorf("HTTP Regular finish"), PoolHTTPRegularFinish, proxyTag)
				resp.Body.Close()
				break
			}
		}
		// Retry
		if resp.StatusCode == http.StatusTooManyRequests {
			ctx.SetPoolContextErrorWithType(fmt.Errorf("errCode:429 errMsg:Acquire proxy failed : no available proxy"), PoolParentProxyFail, proxyTag)
		} else {
			m := make(map[string]interface{})
			err = json.Unmarshal(buf, &m)
			if err != nil {
				Logger.Errorf("proxyHTTPWithConnPool %s Unmarshal resp body failed, body: %s, err: %s", ctx.Req.URL, buf, err)
			} else {
				ctx.SetPoolContextErrorWithType(fmt.Errorf("errCode:%d errMsg:%s", int(m["errCode"].(float64)), m["errMsg"].(string)), PoolParentProxyFail, proxyTag)
			}
		}
		resp.Body.Close()
	}
	if !work {
		// No parentProxyURL works, just return http.StatusTooManyRequests
		Logger.Errorf("proxyHTTPWithConnPool %s cannot find working parent proxy to forward request", ctx.Req.URL.Host)
		rw.WriteHeader(http.StatusTooManyRequests)
		WriteProxyErrorToResponseBody(ctx, rw, http.StatusTooManyRequests, fmt.Sprintf("proxyHTTPWithConnPool %s cannot find working parent proxy to forward request", ctx.Req.URL.Host), tooManyRequests)
		ctx.SetPoolContextErrorWithType(nil, PoolNoAvailableParentProxyFail)
	}
}

func headerContains(header http.Header, name string, value string) bool {
	for _, v := range header[name] {
		for _, s := range strings.Split(v, ",") {
			if strings.EqualFold(value, strings.TrimSpace(s)) {
				return true
			}
		}
	}
	return false
}
