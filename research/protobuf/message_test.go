package protobuf

import (
   "os"
   "testing"
)

func TestMarshal(t *testing.T) {
   checkin := protobuf.Message{
      4: protobuf.Message{ // checkin
         1: protobuf.Message{ // build
            10: protobuf.Varint(29), // sdkVersion
         },
      },
      14: protobuf.Varint(3), // version
      18: protobuf.Message{ // deviceConfiguration
         1: protobuf.Varint(c.Touch_Screen), // touchScreen
         2: protobuf.Varint(c.Keyboard),
         3: protobuf.Varint(c.Navigation),
         4: protobuf.Varint(c.Screen_Layout),
         5: protobuf.Varint(c.Hard_Keyboard),
         6: protobuf.Varint(c.Five_Way_Navigation),
         7: protobuf.Varint(c.Screen_Density),
         8: protobuf.Varint(c.GL_ES_Version),
         11: protobuf.String(platform), // nativePlatform
      },
   }
   for _, library := range c.Shared_Library {
      // .deviceConfiguration.systemSharedLibrary
      checkin.Get(18).Add_String(9, library)
   }
   for _, extension := range c.GL_Extension {
      // .deviceConfiguration.glExtension
      checkin.Get(18).Add_String(15, extension)
   }
   for _, name := range c.Device_Feature {
      // .deviceConfiguration.deviceFeature
      checkin.Get(18).Add(26, protobuf.Message{
         1: protobuf.String(name),
      })
   }
}

func TestUnmarshal(t *testing.T) {
   buf, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   response_wrapper, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   doc_V2 := response_wrapper.Message(1).Message(2).Message(4)
   if v := doc_V2.Message(13).Message(1).Messages(17); len(v) != 4 {
      t.Fatal("File", v)
   }
   if v, err := doc_V2.Message(13).Message(1).Varint(3); err != nil {
      t.Fatal(err)
   } else if v != 10218030 {
      t.Fatal("VersionCode", v)
   }
   if v, err := doc_V2.Message(13).Message(1).String(4); err != nil {
      t.Fatal(err)
   } else if v != "10.21.0" {
      t.Fatal("VersionString", v)
   }
   if v, err := doc_V2.Message(13).Message(1).Varint(9); err != nil {
      t.Fatal(err)
   } else if v != 47705639 {
      t.Fatal("Size", v)
   }
   if v, err := doc_V2.Message(13).Message(1).String(16); err != nil {
      t.Fatal(err)
   } else if v != "Jun 14, 2022" {
      t.Fatal("Date", v)
   }
   if v, err := doc_V2.String(5); err != nil {
      t.Fatal(err)
   } else if v != "Pinterest" {
      t.Fatal("title", v)
   }
   if v, err := doc_V2.String(6); err != nil {
      t.Fatal(err)
   } else if v != "Pinterest" {
      t.Fatal("creator", v)
   }
   if v, err := doc_V2.Message(8).String(2); err != nil {
      t.Fatal(err)
   } else if v != "USD" {
      t.Fatal("currencyCode", v)
   }
   if v, err := doc_V2.Message(13).Message(1).Varint(70); err != nil {
      t.Fatal(err)
   } else if v != 750510010 {
      t.Fatal("NumDownloads", v)
   }
}
