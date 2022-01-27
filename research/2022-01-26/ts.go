package main

import (
   "fmt"
   "net/http"
   "net/url"
)

func main() {
   var req http.Request
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "cbsios-vh.akamaihd.net"
   req.URL.Path = "/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/segment1_1_av.ts"
   req.URL.Scheme = "http"
   val := make(url.Values)
   val["hdntl"] = []string{"exp=1643330565~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=f0923b442f0b0acad2feae9c9bb69f7f9f3d4757f68c56beee82cca989bb6804"}
   val["id"] = []string{"AgBItRcmF8YMPIXq8WFuVofe99NZLY6/oBiPX a7nxe7WVOeZcfLp928J5JEosg7BLeh0j65Jkz2Pw=="}
   req.URL.RawQuery = val.Encode()
   req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"}
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   fmt.Printf("%+v\n", res)
}
