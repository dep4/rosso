package proxychannel

import (
   "crypto/tls"
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
	//CertCache        cert.Cache
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

// WriterWithProtocol .
type WriterWithProtocol struct {
	writer Writer
	length int
}

func (w *WriterWithProtocol) Write(b []byte) (n int, err error) {
	n, err = w.Write(b)
	w.length += n
	return n, err
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

// ReaderWithProtocol .
type ReaderWithProtocol struct {
	reader Reader
	length int
}

func (r *ReaderWithProtocol) Read(b []byte) (n int, err error) {
	n, err = r.Read(b)
	r.length += n
	return n, err
}
