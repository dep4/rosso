package net

import (
   "io"
)

func (v Values) WriteTo3(w io.Writer) (int64, error) {
   var ns int
   for key := range v.Values {
      if n, err := io.WriteString(w, key); err != nil {
         return 0, err
      } else {
         ns += n
      }
      if n, err := io.WriteString(w, "="); err != nil {
         return 0, err
      } else {
         ns += n
      }
      if n, err := io.WriteString(w, v.Get(key)); err != nil {
         return 0, err
      } else {
         ns += n
      }
      if n, err := io.WriteString(w, "\n"); err != nil {
         return 0, err
      } else {
         ns += n
      }
   }
   return int64(ns), nil
}
