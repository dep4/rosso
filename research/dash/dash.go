package main

import (
   "net/http"
   "os"
)

var parts = []string{
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/b7173ce1-3126-4042-a0d1-f8cf8dacd94e/xd9/default_audio128_5_en_main/init0.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/b7173ce1-3126-4042-a0d1-f8cf8dacd94e/xd9/default_audio128_5_en_main/segment0.m4f",
}

func main() {
   file, err := os.Create("ignore/1052529-enc.mp4")
   if err != nil {
      panic(err)
   }
   defer file.Close()
   for _, part := range parts {
      res, err := http.Get(part)
      if err != nil {
         panic(err)
      }
      if res.StatusCode != http.StatusOK {
         panic(err)
      }
      if _, err := file.ReadFrom(res.Body); err != nil {
         panic(err)
      }
      if err := res.Body.Close(); err != nil {
         panic(err)
      }
   }
}
