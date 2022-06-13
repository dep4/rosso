package main

import (
   "fmt"
   "io"
   "net/http"
)

const block = 16 << 7 // 2048

type spy struct {
   io.Reader
}

func (s spy) Read(p []byte) (int, error) {
   n, err := s.Reader.Read(p)
   fmt.Println(n)
   return n, err
}

func main() {
   res, err := http.Get("https://cbsios-vh.akamaihd.net/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/index_0_a.m3u8?null=0&id=AgBItRcmFya85K85p2KGlzWuVQLqG9znXJjcyyrhrcnAQJM+rBxr7wmyn2C8mc2LCajvQqmfeEzkJg%3d%3d&hdntl=exp=1655212847~acl=%2fi%2ftemp_hd_gallery_video%2fCBS_Production_Outlet_VMS%2fvideo_robot%2fCBS_Production_Entertainment%2f2012%2f09%2f12%2f41581439%2fCBS_MELROSE_PLACE_001_SD_prores_78930_*~data=hdntl~hmac=a4af5b6b11debe817bf31f8e86d52701ccb17e1b52cb407b55b100e8c0ada090")
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   fmt.Println(res.ContentLength)
   io.Copy(io.Discard, spy{res.Body})
}
