package net

import (
   "bufio"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strconv"
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
   // GET /fdfe/details?doc=com.instagram.android HTTP/1.1
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

type Values map[string]string

// godocs.io/net/http#Request.Body
func (v Values) Body() io.Reader {
   raw := v.RawQuery()
   return strings.NewReader(raw)
}

// godocs.io/net/http#Request.Header
func (v Values) Header() http.Header {
   vals := make(http.Header)
   for key, val := range v {
      vals.Set(key, val)
   }
   return vals
}

// godocs.io/net/url#URL.RawQuery
func (v Values) RawQuery() string {
   vals := make(url.Values)
   for key, val := range v {
      vals.Set(key, val)
   }
   return vals.Encode()
}
