package net

import (
   "bufio"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strings"
)

func ReadRequest(src io.Reader) (*http.Request, error) {
   text := textproto.NewReader(bufio.NewReader(src))
   line, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   // GET /fdfe/details?doc=com.instagram.android HTTP/1.1
   methodURL := strings.Fields(line)
   if len(methodURL) != 3 {
      return nil, textproto.ProtocolError(line)
   }
   addr, err := url.Parse(methodURL[1])
   if err != nil {
      return nil, err
   }
   addr.Host = head.Get("Host")
   var req http.Request
   req.Body = io.NopCloser(text.R)
   req.Header = http.Header(head)
   req.Method = methodURL[0]
   req.URL = addr
   return &req, nil
}

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ReadQuery(src io.Reader) url.Values {
   vals := make(url.Values)
   buf := bufio.NewReader(src)
   for {
      key, err := buf.ReadString('=')
      if err != nil {
         break
      }
      val, err := buf.ReadString('\n')
      if err != nil {
         break
      }
      vals.Add(key, val)
   }
   return vals
}
