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

func lengthPrefix(bs []byte, vs ...[]byte) []byte {
   var b []byte
   for _, v := range vs {
      b = append(b, v...)
   }
   return protowire.AppendBytes(bs, b)
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
      var b []byte
      b = protowire.AppendTag(b, 1, protowire.BytesType)
      b = lengthPrefix(b,
         protowire.AppendTag(nil, 2, protowire.BytesType),
         protowire.AppendString(nil, "hello"),
      )
      fmt.Printf("%q\n", b)
   }
}
