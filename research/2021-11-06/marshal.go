package main

import (
   "fmt"
   "github.com/segmentio/encoding/proto"
)

var defaultConfig = protoNew{
   1: protoNew{
      1: 1,
      2: 1,
      3: 1,
      4: 1,
      5: true,
      6: true,
      7: 1,
      8: 0x0009_0000,
      10: []string{
         "android.hardware.camera",
         "android.hardware.faketouch",
         "android.hardware.location",
         "android.hardware.screen.portrait",
         "android.hardware.touchscreen",
         "android.hardware.wifi",
      },
      11: []string{
         "armeabi-v7a",
      },
   },
}

type protoOld struct {
   One string `protobuf:"bytes,1"`
   Two string `protobuf:"bytes,2"`
}

type protoNew map[int32]interface{}

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
      fmt.Println(p)
   }
}
