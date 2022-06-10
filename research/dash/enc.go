package main

import (
   "fmt"
   "net/http"
   "os"
)

const init0 = "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/init0.m4f"

var segments = []string{
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment0.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment1.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment2.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment3.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment4.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment5.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment6.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment7.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment8.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/5b025714-8aab-43e4-9ee8-d4f55a14e118/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment9.m4f",
}

func main() {
   enc, err := os.Create("ignore/enc.mp4")
   if err != nil {
      panic(err)
   }
   res, err := http.Get(init0)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   enc.ReadFrom(res.Body)
   for _, segment := range segments {
      fmt.Println(segment)
      res, err := http.Get(segment)
      if err != nil {
         panic(err)
      }
      if res.StatusCode != http.StatusOK {
         panic(res.Status)
      }
      if _, err := enc.ReadFrom(res.Body); err != nil {
         panic(err)
      }
      if err := res.Body.Close(); err != nil {
         panic(err)
      }
   }
}
