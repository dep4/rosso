package main

import (
   "fmt"
   "github.com/segmentio/encoding/proto"
   "google.golang.org/protobuf/encoding/protowire"
)

type protoOld struct {
   One struct {
      Two string `protobuf:"bytes,2"`
   } `protobuf:"bytes,1"`
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
      var p protoOld
      p.One.Two = "hello"
      b, err := proto.Marshal(p)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", b)
   }
   {
      p := protoNew{
         1: protoNew{2: "hello"},
      }
      b := p.marshal()
      fmt.Printf("%q\n", b)
   }
}
