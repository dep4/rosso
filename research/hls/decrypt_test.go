package hls

import (
   "errors"
   "io"
   "net/http"
   "os"
   "testing"
)

func getKey() ([]byte, error) {
   res, err := http.Get("https://cbsios-vh.akamaihd.net/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/crypt.key?null=0&id=AgBItRcmF8YMPMCup2KU42ZeZW8m1CA0O%2fRay%2f9bQ2Hbb1duv9+n8GY8c4ZBp2eIVnJQtUf8GTgCCA%3d%3d&hdntl=exp=1655242816~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_*~data=hdntl~hmac=62caa39a18047ad9a26ca928a6a0b59a4f7256f7bec9536cd208999429019d43")
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   return io.ReadAll(res.Body)
}

func TestHTTP(t *testing.T) {
   res, err := http.Get("https://cbsios-vh.akamaihd.net/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/segment1_0_av.ts?null=0&id=AgBItRcmF8YMPMCup2KU42ZeZW8m1CA0O%2fRay%2f9bQ2Hbb1duv9+n8GY8c4ZBp2eIVnJQtUf8GTgCCA%3d%3d&hdntl=exp=1655242816~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_*~data=hdntl~hmac=62caa39a18047ad9a26ca928a6a0b59a4f7256f7bec9536cd208999429019d43")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      t.Fatal(res.Status)
   }
   key, err := getKey()
   if err != nil {
      t.Fatal(err)
   }
   block, err := newCipher(res.Body, key)
   if err != nil {
      t.Fatal(err)
   }
   dec, err := os.Create("ignore.ts")
   if err != nil {
      t.Fatal(err)
   }
   defer dec.Close()
   if _, err := dec.ReadFrom(block); err != nil {
      t.Fatal(err)
   }
}
