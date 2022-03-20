package main

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "google.golang.org/protobuf/testing/protopack"
)

func main() {
   buf := []byte("Instagram")
   {
      mes, err := protobuf.Unmarshal(buf)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%v %q\n", mes, mes.Marshal())
   }
   {
      var mes protopack.Message
      mes.UnmarshalAbductive(buf, nil)
      fmt.Printf("%v %q\n", mes, mes.Marshal())
   }
}
