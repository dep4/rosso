package http

import (
   "bufio"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strings"
)

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
