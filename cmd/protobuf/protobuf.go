package main

import (
   "bytes"
   "encoding/json"
   "flag"
   "github.com/89z/format/protobuf"
   "os"
)

func main() {
   // f
   var name string
   flag.StringVar(&name, "f", "", "input file")
   // o
   var output string
   flag.StringVar(&output, "o", "", "output file")
   flag.Parse()
   if name != "" {
      buf, err := os.ReadFile(name)
      if err != nil {
         panic(err)
      }
      file, err := os.Create(output)
      if err != nil {
         file = os.Stdout
      }
      defer file.Close()
      mes, err := protobuf.Unmarshal(buf)
      if err != nil {
         panic(err)
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
   } else {
      flag.PrintDefaults()
   }
}
