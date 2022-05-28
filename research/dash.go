package main

import (
   "github.com/edgeware/mp4ff/mp4"
   "net/http"
)

const segment0 = "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/c1cff58f-cc7c-4285-8883-6758d92d8928/xe5/ffe34eec-4812-4ab9-830b-a3f2ea9d076c/segment0.m4f"

func main() {
   res, err := http.Get(segment0)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      panic(res.Status)
   }
   mp4.DecodeFile(res.Body)
}
