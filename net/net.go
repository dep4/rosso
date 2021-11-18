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

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(query []byte) url.Values {
   res := make(url.Values)
   lines := bytes.Split(query, []byte{'\n'})
   for _, line := range lines {
      key, val, ok := keyVal(line)
      if ! ok {
         return nil
      }
      res.Add(string(key), string(val))
   }
   return res
}

func ReadRequest(r io.Reader) (*http.Request, error) {
   t := textproto.NewReader(bufio.NewReader(r))
   s, err := t.ReadLine()
   if err != nil {
      return nil, err
   }
   h, err := t.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   f := strings.Fields(s)
   p, err := url.Parse(f[1])
   if err != nil {
      return nil, err
   }
   p.Host = h.Get("Host")
   return &http.Request{
      Body: io.NopCloser(t.R),
      Header: http.Header(h),
      Method: f[0],
      URL: p,
   }, nil
}

func keyVal(kv []byte) ([]byte, []byte, bool) {
   i := bytes.IndexByte(kv, '=')
   if i == -1 {
      return kv, nil, false
   }
   return kv[:i], kv[i+1:], true
}
