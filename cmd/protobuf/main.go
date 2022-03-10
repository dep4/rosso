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
      err := doProtoBuf()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
