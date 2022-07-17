package main

import (
   "bytes"
   "fmt"
   "github.com/89z/rosso/http"
   "github.com/89z/rosso/strconv"
   "io"
   "net/http/httputil"
   "net/url"
   "os"
   "text/template"
)

func write(req *http.Request, file *os.File) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if file == os.Stdout {
      dump, err := httputil.DumpResponse(res, true)
      if err != nil {
         return err
      }
      var buf string
      if strconv.Valid(dump) {
         buf = string(dump)
      } else {
         buf = strconv.Quote(dump)
      }
      file.WriteString(buf)
   } else {
      buf, err := httputil.DumpResponse(res, false)
      if err != nil {
         return err
      }
      os.Stdout.Write(buf)
      if _, err := file.ReadFrom(res.Body); err != nil {
         return err
      }
   }
   return nil
}

type request_template struct {
   *http.Request
   Body_IO string
   Query url.Values
   Var_Body string
}

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
   temp, err := new(template.Template).Parse(raw_temp)
   if err != nil {
      return err
   }
   return temp.Execute(w, req_temp)
}

const raw_temp = `package main

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
