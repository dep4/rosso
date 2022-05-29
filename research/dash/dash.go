package main

import (
   "net/http"
   "os"
)

var parts = []string{
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/2175d419-9795-4279-9d17-eb6a8dd94ea9/xe3/edfa136e-6e76-4043-bec6-e41e772dff2f/init0.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/2175d419-9795-4279-9d17-eb6a8dd94ea9/xe3/edfa136e-6e76-4043-bec6-e41e772dff2f/segment0.m4f",
}

func main() {
   file, err := os.Create("ignore/1011152-enc.mp4")
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
