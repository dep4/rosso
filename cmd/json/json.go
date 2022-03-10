package main

import (
   "encoding/json"
   "os"
)

func doJSON(input, output string) error {
   src, err := os.Open(input)
   if err != nil {
      return err
   }
   defer src.Close()
   dst, err := os.Create(output)
   if err != nil {
      dst = os.Stdout
   }
   defer dst.Close()
   var val interface{}
   if err := json.NewDecoder(src).Decode(&val); err != nil {
      return err
   }
   enc := json.NewEncoder(dst)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc.Encode(val)
}
