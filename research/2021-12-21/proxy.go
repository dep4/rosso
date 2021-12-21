package main

import (
   "crypto/tls"
   "encoding/json"
   "fmt"
   "github.com/89z/parse/crypto"
   "io"
   "net"
   "net/http"
   "os"
   "os/signal"
   "sync"
   "sync/atomic"
   "time"
   stdcontext "context"
)

func main() {
   // Providing certain log configuration before Run() is optional e.g.
   // ConfigLogging(lconf) where lconf is a *LogConfig
   pc := NewProxychannel(defaultHandlerConfig, defaultServerConfig)
   fmt.Println("runServer")
   pc.runServer()
}

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

// Proxy is a struct that implements ServeHTTP() method
type proxy struct {
   connNum int32
   transport     *http.Transport
   mode          int
}

func newProxy(hconf *http.Transport) *proxy {
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

func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
   fmt.Println(req.Header)
   if req.URL.Host == "" {
      req.URL.Host = req.Host
   }
   atomic.AddInt32(&p.connNum, 1)
   defer func() {
      atomic.AddInt32(&p.connNum, -1)
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
      }
   }
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

// FailEventType. When a request is aborted, the event should be one of the
// following.
const (
   AuthFail           = "AUTH_FAIL"
   BeforeRequestFail  = "BEFORE_REQUEST_FAIL"
   BeforeResponseFail = "BEFORE_RESPONSE_FAIL"
   ConnectFail        = "CONNECT_FAIL"
   HTTPDoRequestFail               = "HTTP_DO_REQUEST_FAIL"
   HTTPRedialCancelTimeout   = "HTTP_REDIAL_CANCEL_TIMEOUT"
   HTTPSRedialCancelTimeout  = "HTTPS_REDIAL_CANCEL_TIMEOUT"
   HTTPWriteClientFail             = "HTTP_WRITE_CLIENT_FAIL"
   ParentProxyFail    = "PARENT_PROXY_FAIL"
   TunnelConnectRemoteFail         = "TUNNEL_CONNECT_REMOTE_FAIL"
   TunnelDialRemoteServerFail      = "TUNNEL_DIAL_REMOTE_SERVER_FAIL"
   TunnelHijackClientConnFail      = "TUNNEL_HIJACK_CLIENT_CONN_FAIL"
   TunnelRedialCancelTimeout = "TUNNEL_REDIAL_CANCEL_TIMEOUT"
   TunnelWriteClientConnFinish     = "TUNNEL_WRITE_CLIENT_CONN_FINISH"
   TunnelWriteEstRespFail          = "TUNNEL_WRITE_EST_RESP_FAIL"
   TunnelWriteTargetConnFinish     = "TUNNEL_WRITE_TARGET_CONN_FINISH"
)

type proxychannel struct {
   server           *http.Server
   waitGroup        *sync.WaitGroup
   serverDone       chan bool
}

func NewProxychannel(hconf *http.Transport, sconf *serverConfig) *proxychannel {
   pc := &proxychannel{
      waitGroup:        &sync.WaitGroup{},
      serverDone:       make(chan bool),
   }
   pc.server = &http.Server{
      Addr:         sconf.ProxyAddr,
      Handler:    newProxy(hconf),
      ReadTimeout:  sconf.ReadTimeout,
      WriteTimeout: sconf.WriteTimeout,
      TLSConfig:    sconf.TLSConfig,
   }
   return pc
}

func (pc *proxychannel) runServer() {
   ctx, cancel := stdcontext.WithCancel(stdcontext.Background())
   defer cancel()
   defer close(pc.serverDone)
   pc.server.BaseContext = func(_ net.Listener) stdcontext.Context { return ctx }
   stop := func() {
   gracefulCtx, cancel := stdcontext.WithTimeout(stdcontext.Background(), 5*time.Second)
   defer cancel()
   if err := pc.server.Shutdown(gracefulCtx); err != nil {
   fmt.Printf("HTTP server Shutdown error: %v\n", err)
   } else {
   fmt.Println("HTTP server gracefully stopped")
   }
   }
   // Run server
   go func() {
   if err := pc.server.ListenAndServe(); err != http.ErrServerClosed {
   //Logger.Errorf("HTTP server ListenAndServe: %v", err)
   os.Exit(1)
   }
   }()
   signalChan := make(chan os.Signal, 1)
   signal.Notify(signalChan)
   // Will block until shutdown signal is received
   <-signalChan
   // Terminate after second signal before callback is done
   go func() {
   <-signalChan
   os.Exit(1)
   }()
   stop()
}

type context struct {
   abort      bool
   data       map[interface{}]interface{}
   err        error
   errType    string
   hijack     bool
   lock       sync.RWMutex
   mitm       bool
   req        *http.Request
   reqLength  int64
   respLength int64
}

func (c *context) setContextErrorWithType(err error, errType string) {
   c.lock.Lock()
   defer c.lock.Unlock()
   if c.errType == HTTPRedialCancelTimeout || c.errType == HTTPSRedialCancelTimeout || c.errType == TunnelRedialCancelTimeout {
      return
   }
   c.errType = errType
   c.err = err
}

func (c *context) setPoolContextErrorWithType(err error, errType string, parentProxy ...string) {
   c.lock.Lock()
   defer c.lock.Unlock()
   switch len(parentProxy) {
   case 0:
      c.errType = errType
      if err != nil {
         if c.err != nil {
            c.err = fmt.Errorf("%s | %s", err, c.err)
         } else {
            c.err = fmt.Errorf("%s", err)
         }
      }
   case 1:
      p := parentProxy[0]
      if err != nil {
         if c.err != nil {
            c.err = fmt.Errorf("(%s) [%s] %s | %s", p, errType, err, c.err)
         } else {
            c.err = fmt.Errorf("(%s) [%s] %s", p, errType, err)
         }
      }
   default:
      return
   }
}

func (c *context) setContextErrType(errType string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.errType == HTTPRedialCancelTimeout || c.errType == HTTPSRedialCancelTimeout || c.errType == TunnelRedialCancelTimeout {
		return
	}
	c.errType = errType
}

var defaultHandlerConfig = &http.Transport{
   TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
   DialContext: (&net.Dialer{
      DualStack: true,
   }).DialContext,
   MaxIdleConns:          100,
   IdleConnTimeout:       90 * time.Second,
   TLSHandshakeTimeout:   10 * time.Second,
   ExpectContinueTimeout: 1 * time.Second,
}

var defaultServerConfig = &serverConfig{
	ProxyAddr:    ":8080",
	ReadTimeout:  60 * time.Second,
	WriteTimeout: 60 * time.Second,
}

type serverConfig struct {
	ProxyAddr    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLSConfig    *tls.Config
}
