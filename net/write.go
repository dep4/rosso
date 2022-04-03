package net

import (
   "bytes"
   "io"
   "net/http"
   "net/url"
   "text/template"
)

const format = `package main

import (
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Body = {{ .BodyIO }}
   req.Header = make(http.Header)
   {{ range $key, $val := .Header -}}
      req.Header[{{ printf "%q" $key }}] = {{ printf "%#v" $val }}
   {{ end -}}
   req.Method = {{ printf "%q" .Method }}
   req.URL = new(url.URL)
   req.URL.Host = {{ printf "%q" .URL.Host }}
   req.URL.Path = {{ printf "%q" .URL.Path }}
   req.URL.RawPath = {{ printf "%q" .URL.RawPath }}
   val := make(url.Values)
   {{ range $key, $val := .Query -}}
      val[{{ printf "%q" $key }}] = {{ printf "%#v" $val }}
   {{ end -}}
   req.URL.RawQuery = val.Encode()
   req.URL.Scheme = {{ printf "%q" .URL.Scheme }}
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
}

var body = strings.NewReader({{ printf "%q" .VarBody }})
`

func WriteRequest(req *http.Request, w io.Writer) error {
   var request requestTemplate
   if req.Body != nil && req.Method != "GET" {
      buf, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      req.Body = io.NopCloser(bytes.NewReader(buf))
      request.BodyIO = "io.NopCloser(body)"
      request.VarBody = string(buf)
   } else {
      request.BodyIO = "io.ReadCloser(nil)"
   }
   request.Query = req.URL.Query()
   request.Request = req
   temp, err := new(template.Template).Parse(format)
   if err != nil {
      return err
   }
   return temp.Execute(w, request)
}

type requestTemplate struct {
   *http.Request
   BodyIO string
   Query url.Values
   VarBody string
}
