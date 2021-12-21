package main

import (
   "crypto/tls"
   "fmt"
   "net"
   "net/http"
   "os"
   "os/signal"
   "sync"
   "time"
   stdcontext "context"
)

func main() {
   // Providing certain log configuration before Run() is optional e.g.
   // ConfigLogging(lconf) where lconf is a *LogConfig
   pc := NewProxychannel(defaultHandlerConfig, defaultServerConfig)
   pc.Run()
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

// Proxychannel is a prxoy server that manages data transmission between http
// clients and destination servers. With the "Extensions" provided by user,
// Proxychannel is able to do authentication, communicate with databases,
// manipulate the requests/responses, etc.
type Proxychannel struct {
	server           *http.Server
	waitGroup        *sync.WaitGroup
	serverDone       chan bool
}

func NewProxychannel(hconf *http.Transport, sconf *serverConfig) *Proxychannel {
	pc := &Proxychannel{
		waitGroup:        &sync.WaitGroup{},
		serverDone:       make(chan bool),
	}
	pc.server = NewServer(hconf, sconf)
	return pc
}

func NewServer(hconf *http.Transport, sconf *serverConfig) *http.Server {
	handler := NewProxy(hconf)
	server := &http.Server{
		Addr:         sconf.ProxyAddr,
		Handler:      handler,
		ReadTimeout:  sconf.ReadTimeout,
		WriteTimeout: sconf.WriteTimeout,
		TLSConfig:    sconf.TLSConfig,
	}
	return server
}

func (pc *Proxychannel) runExtensionManager() {
	defer pc.waitGroup.Done()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	// Will block until shutdown signal is received
	<-signalChan
	// Will block until pc.server has been shut down
	<-pc.serverDone
}

func (pc *Proxychannel) runServer() {
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

// Run launches the ExtensionManager and the HTTP server
func (pc *Proxychannel) Run() {
	pc.waitGroup.Add(1)
	go pc.runExtensionManager()
	pc.runServer()
	pc.waitGroup.Wait()
}


type context struct {
   abort      bool
   closed     bool
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
