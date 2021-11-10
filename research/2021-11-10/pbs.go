package main

import (
   "encoding/json"
   "github.com/89z/parse/protobuf"
   "os"
)

func main() {
   buf, err := os.ReadFile("res.txt")
   if err != nil {
      panic(err)
   }
   dec := protobuf.NewDecoder(buf)
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.SetEscapeHTML(false)
   enc.Encode(dec)
}
