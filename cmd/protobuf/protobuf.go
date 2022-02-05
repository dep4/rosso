package main

import (
   "bytes"
   "encoding/json"
   "flag"
   "fmt"
   "github.com/89z/format/protobuf"
   "os"
)

func main() {
   var output string
   flag.StringVar(&output, "o", "", "output file")
   flag.Parse()
   if flag.NArg() == 1 {
      input := flag.Arg(0)
      bProto, err := os.ReadFile(input)
      if err != nil {
         panic(err)
      }
      file, err := os.Create(output)
      if err != nil {
         file = os.Stdout
      }
      defer file.Close()
      mes, err := protobuf.Unmarshal(bProto)
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
      fmt.Println("protobuf [flags] [file]")
      flag.PrintDefaults()
   }
}
