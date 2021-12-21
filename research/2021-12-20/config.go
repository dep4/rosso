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

type HandlerConfig struct {
	DisableKeepAlive bool
	Delegate         Delegate
	DecryptHTTPS     bool
	Transport        *http.Transport
	Mode             int
}

type Cache struct {
	m sync.Map
}

// Set stores the certificates of hosts that have been seen.
func (c *Cache) Set(host string, cer *tls.Certificate) {
	c.m.Store(host, cer)
}

// Get gets the certificate stored.
func (c *Cache) Get(host string) *tls.Certificate {
	v, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return v.(*tls.Certificate)
}

var DefaultHandlerConfig *HandlerConfig = &HandlerConfig{
	DisableKeepAlive: false,
	Delegate:         &DefaultDelegate{},
	DecryptHTTPS:     false,
	//CertCache:        &Cache{},
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

// DefaultServerConfig .
var DefaultServerConfig *ServerConfig = &ServerConfig{
	ProxyAddr:    ":8080",
	ReadTimeout:  60 * time.Second,
	WriteTimeout: 60 * time.Second,
}

// ServerConfig .
type ServerConfig struct {
	ProxyAddr    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLSConfig    *tls.Config
}

// LogConfig .
type LogConfig struct {
	LoggerName string
	LogLevel   string
	LogOut     string
	LogFormat  string
}

// Writer .
type Writer interface {
	Write([]byte) (int, error)
}

// WriterWithLength .
type WriterWithLength struct {
	writer        interface{} // io.Writer or http.ResponseWriter
	interfaceType int
	length        int
}

func (w *WriterWithLength) Write(b []byte) (n int, err error) {
	if w.interfaceType == 0 {
		// http.ResponseWriter
		respWriter, ok := w.writer.(http.ResponseWriter)
		if !ok {
			panic("w.writer is not a http.ResponseWriter")
		}
		n, err = respWriter.Write(b)
		w.length += n
	} else {
		// io.Writer
		ioWriter, ok := w.writer.(io.Writer)
		if !ok {
			panic("w.writer is not a io.Writer")
		}
		n, err = ioWriter.Write(b)
		w.length += n
	}
	return n, err
}

// Length .
func (w *WriterWithLength) Length() int {
	return w.length
}

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
