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
   req.URL.Path = "/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/crypt.key"
   req.URL.Scheme = "https"
   val := make(url.Values)
   val["hdntl"] = []string{"exp=1643331380~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=07f331ce893a4685a5fda391daf1d2793aad93494ff2c9d0119ca4b4748003ba"}
   val["id"] = []string{"AgBItRcmFy82trTt8WF7NGt/HF0dJn3XUB/YX34dpU2T/4Z6dTL2vKBjrH6VjSYUMovWroYOdvhqSg=="}
   req.URL.RawQuery = val.Encode()
   req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"}
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   fmt.Printf("%+v\n", res)
}
