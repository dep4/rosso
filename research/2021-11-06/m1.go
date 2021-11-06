package main

import (
   "fmt"
)

type protoNew map[int32]interface{}

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

func main() {
   fmt.Println(defaultConfig)
}
