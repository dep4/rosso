package main

import (
   "fmt"
   "github.com/segmentio/encoding/proto"
   "google.golang.org/protobuf/encoding/protowire"
)

type Object struct {
   Array []string `protobuf:"bytes,10"`
}

type (
   array []interface{}
   object map[protowire.Number]interface{}
)

func appendField(out []byte, key protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case string:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendString(out, val)
   case array:
      for _, v := range val {
         out = appendField(out, key, v)
      }
   }
   return out
}

func (o object) marshal() []byte {
   var out []byte
   for key, val := range o {
      out = appendField(out, key, val)
   }
   return out
}

func main() {
   {
      obj := object{
         10: array{"hello", "world"},
      }
      buf := obj.marshal()
      fmt.Printf("%q\n", buf)
   }
   {
      obj := Object{
         Array: []string{"hello", "world"},
      }
      buf, err := proto.Marshal(obj)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", buf)
   }
}
