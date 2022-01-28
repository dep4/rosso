package m3u

import (
   "crypto/cipher"
   "github.com/89z/format"
   "io"
   "net/http"
   "net/url"
   "os"
   "testing"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"

var logLevel format.LogLevel

func getIndex() ([]byte, error) {
   req := new(http.Request)
   req.Header = make(http.Header)
   req.URL = new(url.URL)
   req.URL.Host = "cbsios-vh.akamaihd.net"
   req.URL.Path = "/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/index_1_av.m3u8"
   val := make(url.Values)
   val["hdntl"] = []string{"exp=1643411403~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=7261be21205f1ead3f6936d7b91c014ed22fbba849dbafc068f303ebfc8864cc"}
   val["id"] = []string{"AgBItRcmFy81Sksm82G0S6Vus5DhGvuvBZwDsGQTvpPJN dt XkZKPiuTw6mxQdAIFPdZjWHxM4qug=="}
   req.URL.RawQuery = val.Encode()
   req.URL.Scheme = "https"
   req.Header["User-Agent"] = []string{userAgent}
   logLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

func writeFile(req *http.Request, dec cipher.BlockMode) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   src, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   dst := make([]byte, len(src))
   dec.CryptBlocks(dst, src)
   dst = unpad(dst)
   return os.WriteFile("segment1_1_av.ts", dst, os.ModePerm)
}

func TestDecrypt(t *testing.T) {
   index, err := getIndex()
   if err != nil {
      t.Fatal(err)
   }
   var dec cipher.BlockMode
   for _, form := range Unmarshal(index) {
      req, err := http.NewRequest("GET", form["URI"], nil)
      if err != nil {
         t.Fatal(err)
      }
      req.Header.Set("User-Agent", userAgent)
      if form["METHOD"] != "" {
         dec, err = newDecrypter(req)
         if err != nil {
            t.Fatal(err)
         }
      } else if dec != nil {
         err := writeFile(req, dec)
         if err != nil {
            t.Fatal(err)
         }
         break
      }
   }
}
