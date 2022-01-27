package m3u

import (
   "github.com/89z/format"
   "io"
   "net/http"
   "net/url"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"

var logLevel format.LogLevel

func Index() ([]byte, error) {
   var req http.Request
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "cbsios-vh.akamaihd.net"
   req.URL.Path = "/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/index_1_av.m3u8"
   req.URL.Scheme = "https"
   val := make(url.Values)
   val["hdntl"] = []string{"exp=1643408547~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=c1b9fef5ffe294af9168650f6ffbe577e6b5fc6e466079c6125412c26ec10500"}
   val["id"] = []string{"AgBItRcmFy81SCMb82ELGxzGssIzQl1MceXyxAR0b9Nr ND6lK8zyG456c1r9WmIjYqcm8aSixFxXQ=="}
   req.Header["User-Agent"] = []string{userAgent}
   req.URL.RawQuery = val.Encode()
   res, err := new(http.Transport).RoundTrip(&req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

func segment() ([]byte, error) {
   req := new(http.Request)
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "cbsios-vh.akamaihd.net"
   req.URL.Path = "/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/segment1_1_av.ts"
   req.URL.Scheme = "http"
   val := make(url.Values)
   val["hdntl"] = []string{"exp=1643330565~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=f0923b442f0b0acad2feae9c9bb69f7f9f3d4757f68c56beee82cca989bb6804"}
   val["id"] = []string{"AgBItRcmF8YMPIXq8WFuVofe99NZLY6/oBiPX a7nxe7WVOeZcfLp928J5JEosg7BLeh0j65Jkz2Pw=="}
   req.URL.RawQuery = val.Encode()
   req.Header["User-Agent"] = []string{userAgent}
   logLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

func cryptKey() ([]byte, error) {
   req := new(http.Request)
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "cbsios-vh.akamaihd.net"
   req.URL.Path = "/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/crypt.key"
   req.URL.Scheme = "https"
   val := make(url.Values)
   val["hdntl"] = []string{"exp=1643331380~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=07f331ce893a4685a5fda391daf1d2793aad93494ff2c9d0119ca4b4748003ba"}
   val["id"] = []string{"AgBItRcmFy82trTt8WF7NGt/HF0dJn3XUB/YX34dpU2T/4Z6dTL2vKBjrH6VjSYUMovWroYOdvhqSg=="}
   req.URL.RawQuery = val.Encode()
   req.Header["User-Agent"] = []string{userAgent}
   logLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

