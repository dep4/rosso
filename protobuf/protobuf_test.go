package protobuf

import (
   "fmt"
   "testing"
)

var pinterest = []byte("\n\tPinterest\x18\x88\x80\xbd\x04\"\x069.38.0H\xb6Ƃ\x12R'android.permission.ACCESS_FINE_LOCATIONR'android.permission.ACCESS_NETWORK_STATER$android.permission.ACCESS_WIFI_STATER\x19android.permission.CAMERAR%android.permission.FOREGROUND_SERVICER\x1fandroid.permission.GET_ACCOUNTSR\x1bandroid.permission.INTERNETR\x16android.permission.NFCR android.permission.READ_CONTACTSR(android.permission.READ_EXTERNAL_STORAGER\x1fandroid.permission.READ_PROFILER)android.permission.RECEIVE_BOOT_COMPLETEDR\x1fandroid.permission.RECORD_AUDIOR android.permission.SET_WALLPAPERR\"android.permission.USE_CREDENTIALSR\x1aandroid.permission.VIBRATER\x1candroid.permission.WAKE_LOCKR)android.permission.WRITE_EXTERNAL_STORAGER*com.google.android.c2dm.permission.RECEIVERFcom.google.android.finsky.permission.BIND_GET_INSTALL_REFERRER_SERVICER:com.google.android.providers.gsf.permission.READ_GSERVICESR.com.sec.android.provider.badge.permission.READR/com.sec.android.provider.badge.permission.WRITEZ\x12help@pinterest.comb\x1ahttps://help.pinterest.comj\f500,000,000+r\rcom.pinterestz\x9e\x01Every week we polish up the Pinterest app to make it faster and better than ever. Tell us if you like this newest version at http://help.pinterest.com/contact\x82\x01\vNov 4, 2021\x8a\x01\f\b\x00\x10\x88\x80\xbd\x04\x18\x98\x94\xdb\t\x8a\x01\x18\b\x00\x10\x88\x80\xbd\x04\x18ך7\"\vconfig.ldpi\x8a\x01\x16\b\x00\x10\x88\x80\xbd\x04\x18\x99\x83.\"\tconfig.en\x8a\x01 \b\x00\x10\x88\x80\xbd\x04\x18\xae\x94\xc2\a\"\x12config.armeabi_v7a\xa8\x01\x01\xf2\x01\fContains ads\x92\x02\xb9\x02\x10\xb6Ƃ\x12\x1a#\n\x16com.google.android.gms\x10\xb8\xf9\xf7\x05 \x01(\x000\x01\x1a\x1d\n\x12com.google.ar.core\x10\x90\x88\x8ce(\x010\x02 \x1e(\x03P\x02Z\vconfig.ldpiZ\tconfig.enZ\x12config.armeabi_v7a\x8a\x01\xba\x01\n\xb7\x01AB-xQnrRjA2gOewRVcyJK8QVgDaEb9DrMSv_o5YK27xm4VQ55WGxNAm2hsEYREBVIe5Ja-VbVcukaq65AYBjwBCKHYXO76UTsel3YOHKrEygYroA7MpW5KvRV-s2onK8UkiulWIACDI1aftBYiQXO8UwKI2NF3JO0btX78F9tUfJ3KAFQYMYdfE\x9a\x02\x93\x01\x10\x00\x18\x00\b\x012~\b\x01*xhttps://play-lh.googleusercontent.com/DX5LZpc4SUuvYv2kLA31fi05liLQwizh4DbHYwQthmNP15rul46zG-aquf2Vaf3M9EeBIihWmCsyIWFSLQH\x01:\tPinterest@\x02\x9a\x01\x1btqdNvLiUsPc9jEhccusSR6jwJ8o\x80\x02\x1e\xf2\x02D\x12B\x1a)search?q=pub:Pinterest&o=0&c=3&ksm=1&sb=5P\x03Z\rpub:Pinterest\xf8\x01\x05\xb8\x03\x03\x82\x03\tLifestyle\xa8\x03\x80ʵ\xee\x01\xea\x03\x05500M+\x82\x04\x1c\n\fAug 14, 2012\x12\f\b\x98髁\x05\x10\xc0\x8b\xf6\x82\x02\x92\x04\x00\xa2\x04\v\n\tPinterest\xb0\x04\xd1\xd3¿\x02\xca\x04\b\n\x06\b\aB\x02\b\x03\xd2\x04\b\n\x06\b\aB\x02\b\x0e\xe2\x04\b\n\x06\b\aB\x02\b\x13\xea\x04\x04500M\xf2\x04\v500 million")

var defaultConfig = Object{
   1: Object{
      1: uint64(1),
      2: uint64(1),
      3: uint64(1),
      4: uint64(1),
      5: true,
      6: true,
      7: uint64(1),
      8: uint64(0x0009_0000),
      10: Array{
         "android.hardware.camera",
         "android.hardware.faketouch",
         "android.hardware.location",
         "android.hardware.screen.portrait",
         "android.hardware.touchscreen",
         "android.hardware.wifi",
      },
      11: Array{
         "armeabi-v7a",
      },
   },
}

func TestMarshal(t *testing.T) {
   buf := defaultConfig.Marshal()
   fmt.Printf("%q\n", buf)
   obj := Parse(buf)
   fmt.Println(obj)
}

func TestParse(t *testing.T) {
   fields := Parse(pinterest)
   fmt.Println(fields)
}
