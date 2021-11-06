package main

import (
   "fmt"
   "github.com/segmentio/encoding/proto"
   "google.golang.org/protobuf/encoding/protowire"
)

type Message struct {
   One struct {
      Two string `protobuf:"bytes,2"`
   } `protobuf:"bytes,1"`
}

type message map[protowire.Number]interface{}

func (m message) marshal() []byte {
   var out []byte
   for key, val := range m {
      switch v := val.(type) {
      case message:
         out = protowire.AppendTag(out, key, protowire.BytesType)
         out = protowire.AppendBytes(out, v.marshal())
      case string:
         out = protowire.AppendTag(out, key, protowire.BytesType)
         out = protowire.AppendString(out, v)
      }
   }
   return out
}

func main() {
   {
      var m Message
      m.One.Two = "hello"
      b, err := proto.Marshal(m)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", b)
   }
   {
      m := message{
         1: message{2: "hello"},
      }
      b := m.marshal()
      fmt.Printf("%q\n", b)
   }
}
