package net

import (
   "io"
)

func (v Values) WriteTo1(w io.Writer) (int64, error) {
   var ns int
   write := func(s string) error {
      n, err := io.WriteString(w, s)
      if err != nil {
         return err
      }
      ns += n
      return nil
   }
   for key := range v.Values {
      if err := write(key); err != nil {
         return 0, err
      }
      if err := write("="); err != nil {
         return 0, err
      }
      if err := write(v.Get(key)); err != nil {
         return 0, err
      }
      if err := write("\n"); err != nil {
         return 0, err
      }
   }
   return int64(ns), nil
}
