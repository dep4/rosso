package main

import (
   "crypto/tls"
   "fmt"
   "io"
   "net"
   "net/http"
   "sync"
   "time"
)

type handlerConfig struct {
	DisableKeepAlive bool
	Delegate         Delegate
	DecryptHTTPS     bool
	Transport        *http.Transport
	Mode             int
}

var defaultHandlerConfig = &handlerConfig{
	DisableKeepAlive: false,
	Delegate:         &DefaultDelegate{},
	DecryptHTTPS:     false,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: (&net.Dialer{
			// Timeout:   30 * time.Second,
			// KeepAlive: 30 * time.Second,
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

////////////////////////////////////////////////////////////////////////////////

// Reader .
type Reader interface {
	Read([]byte) (int, error)
}

// Default Settings
const (
	DefaultLoggerName    = "ProxyChannel"
	DefaultLogTimeFormat = "2006-01-02 15:04:05"
	DefaultLogLevel      = "debug"
	DefaultLogOut        = "stderr"
	DefaultLogFormat     = `[%{time:` + DefaultLogTimeFormat + `}] [%{module}] [%{level}] %{message}`
)

// ExtensionManager manage extensions
type ExtensionManager struct {
	extensions map[string]Extension
}

// NewExtensionManager initialize an extension
func NewExtensionManager(m map[string]Extension) *ExtensionManager {
	em := &ExtensionManager{
		extensions: m,
	}
	for ename := range em.extensions {
		em.extensions[ename].SetExtensionManager(em)
	}
	return em
}

// GetExtension get extension by name
func (em *ExtensionManager) GetExtension(name string) (Extension, error) {
   ext, ok := em.extensions[name]
   if !ok {
      return nil, fmt.Errorf("no extension named %s", name)
   }
   return ext, nil
}

// Setup setup all extensions one by one
func (em *ExtensionManager) Setup() {
	var wg sync.WaitGroup
	for name, ext := range em.extensions {
		wg.Add(1)
		go func(name string, ext Extension) {
			defer wg.Done()
			if err := ext.Setup(); err != nil {
				return
			}
		}(name, ext)
	}
	wg.Wait()
}

// Cleanup cleanup all extensions one by one, dont know if the order matters
func (em *ExtensionManager) Cleanup() {
	var wg sync.WaitGroup
	for name, ext := range em.extensions {
		wg.Add(1)
		go func(name string, ext Extension) {
			defer wg.Done()
			if err := ext.Cleanup(); err != nil {
				return
			}
		}(name, ext)
	}
	wg.Wait()
}

// Extension python version __init__(self, engine, **kwargs)
type Extension interface {
	Setup() error
	Cleanup() error
	GetExtensionManager() *ExtensionManager
	SetExtensionManager(*ExtensionManager)
}
