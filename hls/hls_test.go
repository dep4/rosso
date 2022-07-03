package hls

import (
   "fmt"
   "github.com/89z/std/http"
   "io"
   "os"
   "testing"
)

// paramount -b 622520382 -f 499000
const address = "https://cbsios-vh.akamaihd.net/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/index_0_av.m3u8?null=0&id=AgBItRcmFy81SkUfwWIsRdilI6s+0hIRmFI6R378aTEqsuj0TmwsVvPmGEoeaIYYS8H6mKrNRB0PPQ%3d%3d&hdntl=exp=1656910021~acl=%2fi%2ftemp_hd_gallery_video%2fCBS_Production_Outlet_VMS%2fvideo_robot%2fCBS_Production_Entertainment%2f2012%2f09%2f12%2f41581439%2fCBS_MELROSE_PLACE_001_SD_prores_78930_*~data=hdntl~hmac=d571a5878bd4532e7fc553c8a9fd1374e039c9506295dacdcc10533b991a3447"

var client = http.Default_Client

func Test_Block(t *testing.T) {
   res, err := client.Get(address)
   if err != nil {
      t.Fatal(err)
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
   for i, addr := range seg.URI {
      fmt.Println(len(seg.URI)-i)
      res, err := client.Level(0).Get(addr)
      if err != nil {
         t.Fatal(err)
      }
      text, err := io.ReadAll(res.Body)
      if err != nil {
         t.Fatal(err)
      }
      text = block.Decrypt_Key(text)
      if _, err := file.Write(text); err != nil {
         t.Fatal(err)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}

func get_key(s string) ([]byte, error) {
   res, err := client.Get(s)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}

func Test_Stream_Some(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   streams := master.Streams.
      Audio("-ak-").
      Audio("-stereo-").
      Codecs("dvh1")
   for _, stream := range streams {
      fmt.Println(stream)
   }
}

func Test_Stream_All(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
}
