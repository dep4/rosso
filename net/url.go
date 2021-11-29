package net

import (
   "bufio"
   "io"
   "net/url"
   "strings"
)

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(src io.Reader) url.Values {
   vals := make(url.Values)
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      key, val, ok := cutByte(buf.Text(), '=')
      if ok {
         vals.Add(key, val)
      }
   }
   return vals
}

func cutByte(s string, c byte) (string, string, bool) {
   i := strings.IndexByte(s, c)
   if i == -1 {
      return s, "", false
   }
   return s[:i], s[i+1:], true
}
