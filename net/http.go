package net

import (
   "bufio"
   "bytes"
   "fmt"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strings"
   "text/template"
)

func Write_Request(req *http.Request, w io.Writer) error {
   var req_temp request_template
   if req.Body != nil && req.Method != "GET" {
      buf, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      req.Body = io.NopCloser(bytes.NewReader(buf))
      req_temp.Body_IO = "io.NopCloser(body)"
      if bytes.IndexByte(buf, '`') >= 0 {
         req_temp.Var_Body = fmt.Sprintf("%q", buf)
      } else {
         req_temp.Var_Body = fmt.Sprintf("`%s`", buf)
      }
   } else {
      req_temp.Body_IO = "io.ReadCloser(nil)"
   }
   req_temp.Query = req.URL.Query()
   req_temp.Request = req
   temp, err := new(template.Template).Parse(format)
   if err != nil {
      return err
   }
   return temp.Execute(w, req_temp)
}

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
   req.Body = {{ .Body_IO }}
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

var body = strings.NewReader({{ .Var_Body }})
`

type request_template struct {
   *http.Request
   Body_IO string
   Query url.Values
   Var_Body string
}

func Read_Request(in io.Reader) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(bufio.NewReader(in))
   // .Method
   sMethodPath, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   methodPath := strings.Fields(sMethodPath)
   if len(methodPath) != 3 {
      return nil, textproto.ProtocolError(sMethodPath)
   }
   req.Method = methodPath[0]
   // .URL
   addr, err := url.Parse(methodPath[1])
   if err != nil {
      return nil, err
   }
   req.URL = addr
   // .URL.Host
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   if req.URL.Host == "" {
      req.URL.Host = head.Get("Host")
   }
   // .Header
   req.Header = http.Header(head)
   // .Body
   buf := new(bytes.Buffer)
   bLen, err := text.R.WriteTo(buf)
   if err != nil {
      return nil, err
   }
   if bLen >= 1 {
      req.Body = io.NopCloser(buf)
   }
   req.ContentLength = bLen
   return &req, nil
}
