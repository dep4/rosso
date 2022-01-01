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
   "time"
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

func WriteRequest(dst io.Writer, req *http.Request) error {
   now := strconv.FormatInt(time.Now().UnixMicro(), 36)
   if req.Body != nil && req.Method == "POST" {
      buf, err := io.ReadAll(req.Body)
      if err != nil {
         return err
      }
      req.Body = io.NopCloser(bytes.NewReader(buf))
      fmt.Fprintf(dst, "var b%v = strings.NewReader(%q)\n", now, buf)
   }
   fmt.Fprintf(dst, "var r%v = &http.Request{", now)
   // .Method
   fmt.Fprintf(dst, "Method:%q", req.Method)
   // .URL
   fmt.Fprintf(dst, ", URL:%#v", req.URL)
   // .Header
   fmt.Fprintf(dst, ", Header:%#v", req.Header)
   // .Body
   if req.Body != nil && req.Method == "POST" {
      fmt.Fprintf(dst, ", Body:io.NopCloser(b%v)", now)
   }
   fmt.Fprintln(dst, "}")
   return nil
}
