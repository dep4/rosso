package main

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

type (
   array []interface{}
   object map[protowire.Number]interface{}
)

var defaultConfig = object{
   10: array{
      "android.hardware.camera",
      "android.hardware.faketouch",
      "android.hardware.location",
      "android.hardware.screen.portrait",
      "android.hardware.touchscreen",
      "android.hardware.wifi",
   },
}

func main() {
   fmt.Println(defaultConfig)
}
