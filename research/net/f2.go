package net

import (
   "io"
)

func (v Values) WriteTo2(w io.Writer) (int64, error) {
   wr := &writer{w: w}
   for key := range v.Values {
      if _, err := io.WriteString(wr, key); err != nil {
         return 0, err
      }
      if _, err := io.WriteString(wr, "="); err != nil {
         return 0, err
      }
      if _, err := io.WriteString(wr, v.Get(key)); err != nil {
         return 0, err
      }
      if _, err := io.WriteString(wr, "\n"); err != nil {
         return 0, err
      }
   }
   return int64(wr.n), nil
}
