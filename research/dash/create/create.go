package main

import (
   "fmt"
   "net/http"
   "os"
   "path"
)

var segments = []string{
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/init0.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment0.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment1.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment2.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment3.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment4.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment5.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment6.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment7.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment8.m4f",
   "http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/07dd57be-684f-46ac-85e5-d9532d12fc79/xd4/4d629cb1-1728-4a7c-8ad2-8bf439eb7e31/segment9.m4f",
}

func main() {
   for _, segment := range segments {
      fmt.Println(segment)
      res, err := http.Get(segment)
      if err != nil {
         panic(err)
      }
      if res.StatusCode != http.StatusOK {
         panic(res.Status)
      }
      file, err := os.Create("ignore/" + path.Base(segment))
      if err != nil {
         panic(err)
      }
      if _, err := file.ReadFrom(res.Body); err != nil {
         panic(err)
      }
      if err := res.Body.Close(); err != nil {
         panic(err)
      }
      if err := file.Close(); err != nil {
         panic(err)
      }
   }
}
