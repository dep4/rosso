package net

import (
   "bufio"
   "bytes"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strings"
)

func ReadRequest(src io.Reader) (*http.Request, error) {
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
