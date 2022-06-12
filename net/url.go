package net

import (
   "bufio"
   "bytes"
   "io"
   "net/url"
   "strings"
)

type Values struct {
   num int
   r io.Reader
   url.Values
}

func NewValues() *Values {
   var val Values
   val.Values = make(url.Values)
   return &val
}

func (v *Values) Read(p []byte) (int, error) {
   num, err := v.r.Read(p)
   v.num += num
   return num, err
}

func (v Values) ReadFrom(r io.Reader) (int64, error) {
   v.r = r
   scan := bufio.NewScanner(&v)
   for scan.Scan() {
      if scan.Err() != nil {
         return 0, scan.Err()
      }
      key, val, ok := strings.Cut(scan.Text(), "=")
      if ok {
         v.Set(key, val)
      }
   }
   return int64(v.num), nil
}

func (v Values) WriteTo(w io.Writer) (int64, error) {
   var buf bytes.Buffer
   for key := range v.Values {
      buf.WriteString(key)
      buf.WriteByte('=')
      buf.WriteString(v.Get(key))
      buf.WriteByte('\n')
   }
   return buf.WriteTo(w)
}
