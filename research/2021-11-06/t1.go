package main

import (
   "fmt"
   "github.com/89z/parse/protobuf"
   "github.com/segmentio/encoding/proto"
)

type config struct {
   One struct {
      Five bool `protobuf:"varint,5"`
   } `protobuf:"bytes,1"`
}

type object = protobuf.Object

var defaultConfig = object{
   1: object{5: true},
}

func main() {
   {
      var con config
      con.One.Five = true
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
