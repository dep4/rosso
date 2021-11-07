package main

import (
   "fmt"
   "github.com/89z/parse/protobuf"
   "github.com/segmentio/encoding/proto"
)

type config struct {
   One struct {
      Eight int32 `protobuf:"varint,8"`
   } `protobuf:"bytes,1"`
}

type object = protobuf.Object

var defaultConfig = object{
   1: object{
      8: uint64(0x0009_0000),
   },
}

func main() {
   {
      var con config
      con.One.Eight = 0x0009_0000
      buf, err := proto.Marshal(con)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", buf)
   }
   {
      buf := defaultConfig.Marshal()
      fmt.Printf("%q\n", buf)
   }
}
