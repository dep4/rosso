package protobuf

import (
   "fmt"
   "os"
   "testing"
)

// string first:
var files = []string{
   // {4 }:16.49.37 
   "com.google.android.youtube.txt",
   // {4 }:216.1.0.21.137
   "com.instagram.android.txt",
   // {4 }:9.0.12.50
   "com.sec.android.app.launcher.txt",
   // {4 }:G-12.3.10.0
   "com.miui.weather2.txt",
   // {4 }:9.42.0
   "com.pinterest.txt",
   // {4 }:3.4.2
   "org.videolan.vlc.txt",
   // {4 }:5.28.5
   "org.thoughtcrime.securesms.txt",
   // {4 }:2.3.13
   "com.valvesoftware.android.steam.community.txt",
   // {4 }:7.0.704
   "com.xiaomi.smarthome.txt",
   // {4 }:3.53.3
   "com.vimeo.android.videoapp.txt",
   // {4 }:5.8.7
   "com.axis.drawingdesk.v3.txt",
   // {4 }:2.2.4.9
   "com.smarty.voomvoom.txt",
}

func TestProto(t *testing.T) {
   for _, file := range files {
      buf, err := os.ReadFile(file)
      if err != nil {
         t.Fatal(err)
      }
      mes, err := Unmarshal(buf)
      if err != nil {
         t.Fatal(err)
      }
      appDetails := mes.Get(1, 2, 4, 13, 1)
      fmt.Print(file, ":\n", appDetails, "\n")
   }
}
