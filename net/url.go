package net

import (
   "bufio"
   "bytes"
   "io"
   "net/url"
   "strings"
)

type Values struct {
   url.Values
}

func NewValues() Values {
   var val Values
   val.Values = make(url.Values)
   return val
}

func (v Values) ReadFrom(r io.Reader) (int64, error) {
   buf := bufio.NewReader(r)
   var num int
   for {
      key, err := buf.ReadString('=')
      if err != nil {
         return 0, err
      }
      val, err := buf.ReadString('\n')
      num += len(key) + len(val)
      key = strings.TrimSuffix(key, "=")
      val = strings.TrimSuffix(val, "\n")
      v.Set(key, val)
      if err == io.EOF {
         return int64(num), nil
      } else if err != nil {
         return 0, err
      }
   }
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
