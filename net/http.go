package net

import (
   "bytes"
   "bufio"
   "fmt"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strconv"
   "strings"
)

func ReadRequest(src io.Reader, https bool) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(bufio.NewReader(src))
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
   // .URL.Scheme
   if https {
      req.URL.Scheme = "https"
   } else {
      req.URL.Scheme = "http"
   }
   // .URL.Host
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   req.URL.Host = head.Get("Host")
   // .Header
   req.Header = http.Header(head)
   // .ContentLength
   sLength := head.Get("Content-Length")
   if sLength != "" {
      length, err := strconv.ParseInt(sLength, 10, 64)
      if err != nil {
         return nil, err
      }
      req.ContentLength = length
   }
   // .Body
   req.Body = io.NopCloser(text.R)
   return &req, nil
}

func WriteRequest(w io.Writer, q *http.Request) error {
   fmt.Fprintln(w, "package main")
   fmt.Fprintln(w, `import "net/http"`)
   fmt.Fprintln(w, `import "net/http/httputil"`)
   fmt.Fprintln(w, `import "net/url"`)
   fmt.Fprintln(w, `import "os"`)
   if q.Body != nil && q.Method == "POST" {
      buf, err := io.ReadAll(q.Body)
      if err != nil {
         return err
      }
      q.Body = io.NopCloser(bytes.NewReader(buf))
      fmt.Fprintln(w, `import "io"`)
      fmt.Fprintln(w, `import "strings"`)
      fmt.Fprintf(w, "var body = strings.NewReader(%q)\n", buf)
   }
   fmt.Fprintln(w, "func main() {")
   fmt.Fprintln(w, "var q http.Request")
   if q.Body != nil && q.Method == "POST" {
      fmt.Fprintln(w, "q.Body=io.NopCloser(body)")
   }
   fmt.Fprintf(w, "q.Method=%q\n", q.Method)
   fmt.Fprintln(w, "q.URL=new(url.URL)")
   if q.URL.RawQuery != "" {
      val, err := url.ParseQuery(q.URL.RawQuery)
      if err != nil {
         return err
      }
      fmt.Fprintf(w, "q.URL.RawQuery=%#v.Encode()\n", val)
   }
   fmt.Fprintf(w, "q.URL.Host=%q\n", q.URL.Host)
   fmt.Fprintf(w, "q.URL.Path=%q\n", q.URL.Path)
   fmt.Fprintf(w, "q.URL.Scheme=%q\n", q.URL.Scheme)
   fmt.Fprintln(w, "q.Header=make(http.Header)")
   for key, val := range q.Header {
      fmt.Fprintf(w, "q.Header[%q]=%#v\n", key, val)
   }
   fmt.Fprintln(w, "s, err := new(http.Transport).RoundTrip(&q)")
   fmt.Fprintln(w, "if err != nil { panic(err) }")
   fmt.Fprintln(w, "defer s.Body.Close()")
   fmt.Fprintln(w, "buf, err := httputil.DumpResponse(s, true)")
   fmt.Fprintln(w, "if err != nil { panic(err) }")
   fmt.Fprintln(w, "os.Stdout.Write(buf)}")
   return nil
}
