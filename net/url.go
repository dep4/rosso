package net

import (
   "bufio"
   "github.com/89z/parse/strings"
   "io"
   "net/url"
)

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(src io.Reader) url.Values {
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
