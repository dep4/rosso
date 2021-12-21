package proxychannel

import (
   "fmt"
   "github.com/op/go-logging"
   "io"
   "io/ioutil"
   "os"
   "strings"
   "sync"
)

var rootLoggerName string = "ProxyChannel"

// Logger is used to print log in proxychannel
var Logger *logging.Logger = logging.MustGetLogger(rootLoggerName)

// Default Settings
const (
	DefaultLoggerName    = "ProxyChannel"
	DefaultLogTimeFormat = "2006-01-02 15:04:05"
	DefaultLogLevel      = "debug"
	DefaultLogOut        = "stderr"
	DefaultLogFormat     = `[%{time:` + DefaultLogTimeFormat + `}] [%{module}] [%{level}] %{message}`
)

func init() {
	var out io.Writer

	switch DefaultLogOut {
	case "stderr":
		out = os.Stderr
	case "stdout":
		out = os.Stdout
	default:
		out = ioutil.Discard
	}

	backend := logging.NewLogBackend(out, "", 0)
	logging.SetBackend(backend)

	l := logging.GetLevel(DefaultLogLevel)
	logging.SetLevel(l, DefaultLoggerName)

	formatter := logging.MustStringFormatter(DefaultLogFormat)
	logging.SetFormatter(formatter)
}

// ConfigLogging sets the log style.
func ConfigLogging(conf *LogConfig) error {
	if err := SetLoggingBackend(conf.LogOut); err != nil {
		return err
	}
	if err := SetLoggingFormat(conf.LogFormat); err != nil {
		return err
	}
	debug := false
	if conf.LogLevel == "debug" {
		debug = true
	}
	if err := SetLoggingLevel(conf.LogLevel, debug); err != nil {
		return err
	}

	return nil
}

// SetLoggingLevel .
func SetLoggingLevel(level string, debug bool) error {

	if strings.TrimSpace(level) == "" {
		level = DefaultLogLevel
	}
	var logLevel logging.Level
	var err error
	if logLevel, err = logging.LogLevel(level); err != nil {
		return err
	}

	if debug {
		logLevel = logging.DEBUG
	}
	logging.SetLevel(logLevel, DefaultLoggerName)
	return nil
}

// SetLoggingFormat .
func SetLoggingFormat(format string) error {
	var formatter logging.Formatter
	var err error
	if formatter, err = logging.NewStringFormatter(format); err != nil {
		return err
	}
	logging.SetFormatter(formatter)
	return nil
}

// SetLoggingBackend .
func SetLoggingBackend(out string) error {
	var o io.Writer
	switch out {
	case "stdout":
		o = os.Stdout
	case "stderr", "":
		o = os.Stderr
	default:
		f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

		if err != nil {
			return err
		}

		o = f
	}

	backend := logging.NewLogBackend(o, "", 0)
	logging.SetBackend(backend)
	return nil
}

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
		return nil, fmt.Errorf("No extension named %s", name)
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
			Logger.Infof("Extension [%s] Setup start!\n", name)
			if err := ext.Setup(); err != nil {
				Logger.Errorf("Extension [%s] Setup error: %v\n", name, err)
				return
			}
			Logger.Infof("Extension [%s] Setup done!\n", name)
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
			Logger.Infof("Extension [%s] Cleanup start!\n", name)
			if err := ext.Cleanup(); err != nil {
				Logger.Errorf("Extension [%s] Cleanup error: %v\n", name, err)
				return
			}
			Logger.Infof("Extension [%s] Cleanup done!\n", name)
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
