package net

import (
   "bufio"
   "github.com/89z/parse/strings"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   stdstr "strings"
)

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ReadQuery(src io.Reader) url.Values {
   vals := make(url.Values)
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      key, val, ok := strings.CutByte(buf.Text(), '=')
      if ok {
         vals.Add(key, val)
      }
   }
   return vals
}

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
   methodURL := stdstr.Fields(line)
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
