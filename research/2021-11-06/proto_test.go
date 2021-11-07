package protobuf

import (
   "fmt"
   "testing"
)

type badArray array

var defaultConfig = object{
   1: object{
      10: badArray{
         "android.hardware.camera",
         "android.hardware.faketouch",
         "android.hardware.location",
         "android.hardware.screen.portrait",
         "android.hardware.touchscreen",
         "android.hardware.wifi",
      },
   },
}

func TestProto(t *testing.T) {
   buf := defaultConfig.marshal()
   fmt.Printf("%q\n", buf)
}
