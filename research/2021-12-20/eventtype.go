package main

import (
   "context"
   "fmt"
   "net"
   "net/http"
   "os"
   "os/signal"
   "sync"
   "time"
)

func main() {
   // Providing certain log configuration before Run() is optional e.g.
   // ConfigLogging(lconf) where lconf is a *LogConfig
   pc := NewProxychannel(
      defaultHandlerConfig, defaultServerConfig, make(map[string]Extension),
   )
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
	extensionManager *ExtensionManager
	server           *http.Server
	waitGroup        *sync.WaitGroup
	serverDone       chan bool
}

func NewProxychannel(hconf *handlerConfig, sconf *serverConfig, m map[string]Extension) *Proxychannel {
	pc := &Proxychannel{
		extensionManager: NewExtensionManager(m),
		waitGroup:        &sync.WaitGroup{},
		serverDone:       make(chan bool),
	}
	pc.server = NewServer(hconf, sconf, pc.extensionManager)
	return pc
}

func NewServer(hconf *handlerConfig, sconf *serverConfig, em *ExtensionManager) *http.Server {
	handler := NewProxy(hconf, em)
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
	go pc.extensionManager.Setup() // TODO: modify setup and error handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	// Will block until shutdown signal is received
	<-signalChan
	// Will block until pc.server has been shut down
	<-pc.serverDone
	pc.extensionManager.Cleanup()
}

func (pc *Proxychannel) runServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer close(pc.serverDone)
	pc.server.BaseContext = func(_ net.Listener) context.Context { return ctx }
	stop := func() {
		gracefulCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
