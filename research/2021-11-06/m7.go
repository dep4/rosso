package main

import (
   "fmt"
   "github.com/segmentio/encoding/proto"
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(out []byte, key protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case array:
      for _, v := range val {
         out = appendField(out, key, v)
      }
   case bool:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, protowire.EncodeBool(val))
   case object:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendBytes(out, val.marshal())
   case string:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendString(out, val)
   case uint64:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, val)
   }
   return out
}

type array []interface{}

type object map[protowire.Number]interface{}

func (o object) marshal() []byte {
   var out []byte
   for key, val := range o {
      out = appendField(out, key, val)
   }
   return out
}

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

type Config struct {
   DeviceConfiguration DeviceConfiguration `protobuf:"bytes,1"`
}

type DeviceConfiguration struct {
   TouchScreen            int32   `protobuf:"varint,1"`
   Keyboard               int32   `protobuf:"varint,2"`
   Navigation             int32   `protobuf:"varint,3"`
   ScreenLayout           int32   `protobuf:"varint,4"`
   HasHardKeyboard        bool    `protobuf:"varint,5"`
   HasFiveWayNavigation   bool    `protobuf:"varint,6"`
   ScreenDensity          int32   `protobuf:"varint,7"`
   GlEsVersion            int32   `protobuf:"varint,8"`
   SystemAvailableFeature []string `protobuf:"bytes,10"`
   NativePlatform         []string `protobuf:"bytes,11"`
}

var DefaultConfig = Config{
   DeviceConfiguration{
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

func main() {
   {
      buf, err := proto.Marshal(DefaultConfig)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", buf)
   }
   {
      buf := defaultConfig.marshal()
      fmt.Printf("%q\n", buf)
   }
}
