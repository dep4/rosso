package main

import (
   "fmt"
   "github.com/segmentio/encoding/proto"
   "google.golang.org/protobuf/encoding/protowire"
)

type protoOld struct {
   One string `protobuf:"bytes,1"`
   Two string `protobuf:"bytes,2"`
}

type protoNew map[protowire.Number]interface{}

func (p protoNew) marshal() []byte {
   var b []byte
   for k, v := range p {
      s, ok := v.(string)
      if ok {
         b = protowire.AppendTag(b, k, protowire.BytesType)
         b = protowire.AppendString(b, s)
      }
   }
   return b
}

func main() {
   {
      p := protoOld{One: "hello", Two: "world"}
      b, err := proto.Marshal(p)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", b)
   }
   {
      p := protoNew{1: "hello", 2: "world"}
      b := p.marshal()
      fmt.Printf("%q\n", b)
   }
}
