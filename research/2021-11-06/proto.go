package main

import (
   "fmt"
   "github.com/89z/parse/protobuf"
)

type (
   array protobuf.Array
   object = protobuf.Object
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
   buf := defaultConfig.Marshal()
   fmt.Printf("%q\n", buf)
}
