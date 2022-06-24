package http

import (
   "bufio"
   "bytes"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strings"
)

var NewRequest = http.NewRequest

type (
   Request = http.Request
   Transport = http.Transport
)

func Read_Request(in io.Reader) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(bufio.NewReader(in))
   // .Method
   raw_method_path, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   method_path := strings.Fields(raw_method_path)
   req.Method = method_path[0]
   // .URL
   addr, err := url.Parse(method_path[1])
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
   length, err := text.R.WriteTo(buf)
   if err != nil {
      return nil, err
   }
   if length >= 1 {
      req.Body = io.NopCloser(buf)
   }
   req.ContentLength = length
   return &req, nil
}
