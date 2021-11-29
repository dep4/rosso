package net

import (
   "bufio"
   "github.com/89z/parse/strings"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
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
   method, sURL, ok := strings.CutByte(line, ' ')
   if !ok {
      return nil, textproto.ProtocolError(line)
   }
   tURL, err := url.Parse(sURL)
   if err != nil {
      return nil, err
   }
   tURL.Host = head.Get("Host")
   var req http.Request
   req.Body = io.NopCloser(text.R)
   req.Header = http.Header(head)
   req.Method = method
   req.URL = tURL
   return &req, nil
}
