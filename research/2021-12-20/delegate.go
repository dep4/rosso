package main

import (
   "crypto/tls"
   "fmt"
   "net"
   "net/http"
   "sync"
   "time"
)

type context struct {
	Req        *http.Request
	Data       map[interface{}]interface{}
	abort      bool
	Hijack     bool
	MITM       bool
	ReqLength  int64
	RespLength int64
	ErrType    string
	Err        error
	Closed     bool
	Lock       sync.RWMutex
}

func (c *context) setContextErrorWithType(err error, errType string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if c.ErrType == HTTPRedialCancelTimeout || c.ErrType == HTTPSRedialCancelTimeout || c.ErrType == TunnelRedialCancelTimeout {
		return
	}
	c.ErrType = errType
	c.Err = err
}

////////////////////////////////////////////////////////////////////////////////

func (c *context) setPoolContextErrorWithType(err error, errType string, parentProxy ...string) {
   c.Lock.Lock()
   defer c.Lock.Unlock()
   switch len(parentProxy) {
   case 0:
   c.ErrType = errType
   if err != nil {
   if c.Err != nil {
   c.Err = fmt.Errorf("%s | %s", err, c.Err)
   } else {
   c.Err = fmt.Errorf("%s", err)
   }
   }
   case 1:
   p := parentProxy[0]
   if err != nil {
   if c.Err != nil {
   c.Err = fmt.Errorf("(%s) [%s] %s | %s", p, errType, err, c.Err)
   } else {
   c.Err = fmt.Errorf("(%s) [%s] %s", p, errType, err)
   }
   }
   default:
   return
   }
}

// SetContextErrType .
func (c *context) SetContextErrType(errType string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if c.ErrType == HTTPRedialCancelTimeout || c.ErrType == HTTPSRedialCancelTimeout || c.ErrType == TunnelRedialCancelTimeout {
		return
	}
	c.ErrType = errType
}

// SetContextError .
func (c *context) SetContextError(err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.Err = err
}

// Abort sets abort to true.
func (c *context) Abort() {
	c.abort = true
}

// AbortWithError sets Err and abort to true.
func (c *context) AbortWithError(err error) {
	c.Lock.Lock()
	c.Err = err
	c.Lock.Unlock()
	c.abort = true
}

// IsAborted checks whether abort is set to true.
func (c *context) IsAborted() bool {
	return c.abort
}

type handlerConfig struct {
	DisableKeepAlive bool
	//Delegate         Delegate
	DecryptHTTPS     bool
	Transport        *http.Transport
	Mode             int
}

var defaultHandlerConfig = &handlerConfig{
	DisableKeepAlive: false,
	//Delegate:         &DefaultDelegate{},
	DecryptHTTPS:     false,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: (&net.Dialer{
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
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
