package main

import (
   "bytes"
   "errors"
   "flag"
   "fmt"
   "github.com/agatan/groxy"
   "io/ioutil"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "sort"
   "strings"
   "unicode/utf8"
   shellquote "github.com/kballard/go-shellquote"
)

var ErrNonUTF8Body = errors.New("request body is not utf-8 string")

type Command []string

func Curlize(r *http.Request) (Command, error) {
	var command []string

	command = append(command, "curl", "-X", r.Method)

	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if err := r.Body.Close(); err != nil {
			return nil, err
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		if !utf8.Valid(body) {
			return nil, ErrNonUTF8Body
		}
		if len(body) != 0 {
			command = append(command, "-d", string(body))
		}
	}

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command = append(command, "-H", fmt.Sprintf("%s: %s", k, strings.Join(r.Header[k], " ")))
	}

	command = append(command, r.URL.String())

	return command, nil
}

func (c Command) String() string {
	return shellquote.Join(c...)
}


////////////////////////////////////////////////////////////////////////////////
// main

type options struct {
	address string
	reverse string
}

func main() {
	var op options

	flag.StringVar(&op.address, "addr", ":8080", "listening address")
	flag.StringVar(&op.reverse, "reverse", "", "upstream address for reverse proxy mode")

	flag.Parse()

	proxy := &groxy.ProxyServer{
		HTTPSAction: groxy.HTTPSActionMITM,
	}
	proxy.Use(func(h groxy.Handler) groxy.Handler {
		return func(r *http.Request) (*http.Response, error) {
			curl, err := Curlize(r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println(curl.String())
			}
			return h(r)
		}
	})

	if op.reverse != "" {
		upstream, err := url.Parse(op.reverse)
		if err != nil {
			panic(err)
		}
		proxy.NonProxyRequestHandler = reverseProxyHandler(upstream)
	}

	if err := http.ListenAndServe(op.address, proxy); err != nil {
		panic(err)
	}
}

func reverseProxyHandler(upstream *url.URL) http.Handler {
	reverseProxy := httputil.NewSingleHostReverseProxy(upstream)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := *upstream
		u.Path = r.URL.Path
		u.RawQuery = r.URL.Query().Encode()
		r.URL = &u
		r.Host = u.Host
		curl, err := Curlize(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println(curl.String())
		}
		reverseProxy.ServeHTTP(w, r)
	})
}
