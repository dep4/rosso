package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/agatan/curlize"
	"github.com/agatan/groxy"
)

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
			curl, err := curlize.Curlize(r)
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
		curl, err := curlize.Curlize(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println(curl.String())
		}
		reverseProxy.ServeHTTP(w, r)
	})
}
