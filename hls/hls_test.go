package hls

import (
   "errors"
   "fmt"
   "io"
   "net/http"
   "os"
   "testing"
)

func Test_Segment(t *testing.T) {
   file, err := os.Open("ignore/apple-audio.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   seg, err := New_Scanner(file).Segment()
   if err != nil {
      t.Fatal(err)
   }
   for _, pro := range seg.Protected {
      fmt.Println(pro)
   }
   fmt.Println(seg.Key)
}

func Test_Block(t *testing.T) {
   res, err := http.Get("https://cbsios-vh.akamaihd.net/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/index_0_av.m3u8?null=0&id=AgBItRcmF8YMPETJp2Idb%2ff8kST9HgI7mEbBnb7XI96bqUv7h7HvAzf5egQq8EdGCZGfDgozAsOiGw%3d%3d&hdntl=exp=1655249604~acl=%2fi%2ftemp_hd_gallery_video%2fCBS_Production_Outlet_VMS%2fvideo_robot%2fCBS_Production_Entertainment%2f2012%2f09%2f12%2f41581439%2fCBS_MELROSE_PLACE_001_SD_prores_78930_*~data=hdntl~hmac=9e7582fede5fb810be51146b848d2df4e675ed8d78d39931da3273f5880dcfa2")
   if err != nil {
      t.Fatal(err)
   }
   if res.StatusCode != http.StatusOK {
      t.Fatal(res.Status)
   }
   seg, err := New_Scanner(res.Body).Segment()
   if err != nil {
      t.Fatal(err)
   }
   if err := res.Body.Close(); err != nil {
      t.Fatal(err)
   }
   key, err := get_key(seg.Key)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.ts")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   block, err := New_Block(key)
   if err != nil {
      t.Fatal(err)
   }
   for i, addr := range seg.Protected {
      fmt.Println(len(seg.Protected)-i)
      res, err := http.Get(addr)
      if err != nil {
         t.Fatal(err)
      }
      if res.StatusCode != http.StatusOK {
         t.Fatal(res.Status)
      }
      if _, err := file.ReadFrom(block.Mode_Key(res.Body)); err != nil {
         t.Fatal(err)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}

func get_key(s string) ([]byte, error) {
   res, err := http.Get(s)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   return io.ReadAll(res.Body)
}

var raw_ivs = []string{
   "00000000000000000000000000000001",
   "0X00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func Test_Hex(t *testing.T) {
   for _, raw_iv := range raw_ivs {
      iv, err := Segment{Raw_IV: raw_iv}.IV()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(iv)
   }
}
