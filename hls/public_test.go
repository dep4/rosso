package hls

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "testing"
)

func TestMaster(t *testing.T) {
   file, err := os.Open("m3u8/roku-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := NewScanner(file).Master(&url.URL{})
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
   for _, media := range master.Media {
      fmt.Println(media)
   }
   fmt.Println(master.Stream(7996499))
}

func TestSegment(t *testing.T) {
   seg, err := newSegment()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println("GET", seg.Key)
   res, err := http.Get(seg.Key.String())
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   block, err := NewCipher(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for _, info := range seg.Info {
      fmt.Println("GET", info.URI)
      res, err := http.Get(info.URI.String())
      if err != nil {
         t.Fatal(err)
      }
      if res.StatusCode != http.StatusOK {
         t.Fatal(res.Status)
      }
      if _, err := block.Copy(file, res.Body, info.IV); err != nil {
         t.Fatal(err)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}

const cbcSegment = 
   "https://cbcrcott-gem.akamaized.net/0f73fb9d-87f0-4577-81d1-e6e970b89a69" +
   "/CBC_DOWNTON_ABBEY_S01E05.ism/QualityLevels(400044)" +
   "/Manifest(video,format=m3u8-aapl,filter=desktop)"

func newSegment() (*Segment, error) {
   file, err := os.Open("m3u8/cbc-video.m3u8")
   if err != nil {
      return nil, err
   }
   defer file.Close()
   addr, err := url.Parse(cbcSegment)
   if err != nil {
      return nil, err
   }
   return NewScanner(file).Segment(addr)
}
