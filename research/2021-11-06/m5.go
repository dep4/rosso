package main

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

func (o object) marshal() []byte {
   var out []byte
   for key, val := range o {
      switch v := val.(type) {
      case object:
         out = protowire.AppendTag(out, key, protowire.BytesType)
         out = protowire.AppendBytes(out, v.marshal())
      case string:
         out = protowire.AppendTag(out, key, protowire.BytesType)
         out = protowire.AppendString(out, v)
      case bool:
         out = protowire.AppendTag(out, key, protowire.VarintType)
         out = protowire.AppendVarint(out, protowire.EncodeBool(v))
      case uint64:
         out = protowire.AppendTag(out, key, protowire.VarintType)
         out = protowire.AppendVarint(out, v)
      }
   }
   return out
}

type (
   array []interface{}
   object map[protowire.Number]interface{}
)

var defaultConfig = object{
   1: object{
      1: uint64(1),
      2: uint64(1),
      3: uint64(1),
      4: uint64(1),
      5: true,
      6: true,
      7: uint64(1),
      8: uint64(0x0009_0000),
      10: array{
         "android.hardware.camera",
         "android.hardware.faketouch",
         "android.hardware.location",
         "android.hardware.screen.portrait",
         "android.hardware.touchscreen",
         "android.hardware.wifi",
      },
      11: array{
         "armeabi-v7a",
      },
   },
}

func main() {
   buf := defaultConfig.marshal()
   fmt.Printf("%q\n", buf)
}
