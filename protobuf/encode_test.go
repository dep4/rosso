package protobuf

import (
   "fmt"
   "testing"
)

type uploadDeviceConfigRequest struct {
   DeviceConfiguration deviceConfigurationProto `json:"1"`
}

type deviceConfigurationProto struct {
   TouchScreen int32 `json:"1"`
   Keyboard int32 `json:"2"`
   Navigation int32 `json:"3"`
   ScreenLayout int32 `json:"4"`
   HasHardKeyboard bool `json:"5"`
   HasFiveWayNavigation bool `json:"6"`
   ScreenDensity int32 `json:"7"`
   GlEsVersion int32 `json:"8"`
   SystemAvailableFeature []string `json:"10"`
   NativePlatform []string `json:"11"`
}

func TestEncode(t *testing.T) {
   defaultConfig := uploadDeviceConfigRequest{
      DeviceConfiguration: deviceConfigurationProto{
         TouchScreen: 1,
         Keyboard: 1,
         Navigation: 1,
         ScreenLayout: 1,
         HasHardKeyboard: true,
         HasFiveWayNavigation: true,
         ScreenDensity: 1,
         GlEsVersion: 0x0009_0000,
         SystemAvailableFeature: []string{
            "android.hardware.camera",
            "android.hardware.faketouch",
            "android.hardware.location",
            "android.hardware.screen.portrait",
            "android.hardware.touchscreen",
            "android.hardware.wifi",
         },
         NativePlatform: []string{
            "armeabi-v7a",
         },
      },
   }
   enc, err := NewEncoder(defaultConfig)
   if err != nil {
      t.Fatal(err)
   }
   buf, err := enc.Encode()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", buf)
}
