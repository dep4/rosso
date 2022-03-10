package main

import (
   "bytes"
   "encoding/json"
   "flag"
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
   indent := new(bytes.Buffer)
   enc := json.NewEncoder(indent)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   if err := enc.Encode(mes); err != nil {
      panic(err)
   }
   if _, err := file.ReadFrom(indent); err != nil {
      panic(err)
   }
}
