package net

import (
   "bufio"
   "io"
   "net/url"
   "strings"
)

type Values struct {
   url.Values
}

func New_Values() Values {
   var val Values
   val.Values = make(url.Values)
   return val
}

func (v Values) ReadFrom(r io.Reader) (int64, error) {
   buf := bufio.NewReader(r)
   var num int
   for {
      key, err := buf.ReadString('=')
      if err == io.EOF {
         break
      } else if err != nil {
         return 0, err
      }
      val, err := buf.ReadString('\n')
      num += len(key) + len(val)
      key = strings.TrimSuffix(key, "=")
      val = strings.TrimSuffix(val, "\n")
      v.Set(key, val)
      if err == io.EOF {
         break
      } else if err != nil {
         return 0, err
      }
   }
   return int64(num), nil
}

func (v Values) WriteTo(w io.Writer) (int64, error) {
   var ns int64
   write := func(ss ...string) error {
      for _, s := range ss {
         n, err := io.WriteString(w, s)
         if err != nil {
            return err
         }
         ns += int64(n)
      }
      return nil
   }
   for key := range v.Values {
      err := write(key, "=", v.Get(key), "\n")
      if err != nil {
         return 0, err
      }
   }
   return ns, nil
}
