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
   closed     bool
   data       map[interface{}]interface{}
   errType    string
   hijack     bool
   respLength int64
   reqLength  int64
   req        *http.Request
   lock       sync.RWMutex
   abort      bool
   mitm       bool
   
   err        error
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

type handlerConfig struct {
	DisableKeepAlive bool
	DecryptHTTPS     bool
	Transport        *http.Transport
	Mode             int
}

var defaultHandlerConfig = &handlerConfig{
	DisableKeepAlive: false,
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
