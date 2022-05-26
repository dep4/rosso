package main

import (
   "net/http"
   "os"
)

var segments = []string{
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6240731308001/edb20cb4-b298-4f36-b8ee-38ba53cfd0b9/xa9/87111d42-1b26-40a3-9270-8c6212d0a818/init.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6240731308001/cc9e43e2-d7c2-44ff-a56e-560584a31527/xa9/87111d42-1b26-40a3-9270-8c6212d0a818/segment0.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6240731308001/cc9e43e2-d7c2-44ff-a56e-560584a31527/xa9/87111d42-1b26-40a3-9270-8c6212d0a818/segment1.m4f",
}

func main() {
   file, err := os.Create("ignore/enc.mp4")
   if err != nil {
      panic(err)
   }
   defer file.Close()
   for _, segment := range segments {
      res, err := http.Get(segment)
      if err != nil {
         panic(err)
      }
      if res.StatusCode != http.StatusOK {
         panic(res.Status)
      }
      if _, err := file.ReadFrom(res.Body); err != nil {
         panic(err)
      }
      if err := res.Body.Close(); err != nil {
         panic(err)
      }
   }
}
