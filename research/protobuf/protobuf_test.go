package protobuf

import (
   "fmt"
   "testing"
)

func Test_Add(t *testing.T) {
   androids := []string{
      "android.hardware.bluetooth",
      "android.hardware.bluetooth_le",
      "android.hardware.camera",
      "android.hardware.camera.autofocus",
      "android.hardware.camera.front",
      "android.hardware.location",
      "android.hardware.location.gps",
      "android.hardware.location.network",
      "android.hardware.microphone",
      "android.hardware.opengles.aep",
      "android.hardware.screen.landscape",
      "android.hardware.screen.portrait",
      "android.hardware.sensor.accelerometer",
      "android.hardware.sensor.compass",
      "android.hardware.sensor.gyroscope",
      "android.hardware.telephony",
      "android.hardware.touchscreen",
      "android.hardware.touchscreen.multitouch",
      "android.hardware.usb.host",
      "android.hardware.wifi",
      "android.software.device_admin",
      "android.software.midi",
   }
   checkin := make(Message)
   for range androids {
      err := checkin.Add(1, Varint(2))
      if err != nil {
         t.Fatal(err)
      }
   }
   fmt.Println(checkin)
}
