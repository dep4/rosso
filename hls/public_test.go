package hls

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "sort"
   "testing"
)

const cbcMaster =
   "https://cbcrcott-gem.akamaized.net/0f73fb9d-87f0-4577-81d1-e6e970b89a69" +
   "/CBC_DOWNTON_ABBEY_S01E05.ism/desktop_master.m3u8"

func TestMaster(t *testing.T) {
   file, err := os.Open("m3u8/cbc-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   addr, err := url.Parse(cbcMaster)
   if err != nil {
      t.Fatal(err)
   }
   master, err := NewScanner(file).Master(addr)
   if err != nil {
      t.Fatal(err)
   }
   for i, video := range master.Stream {
      if i == 0 {
         audio := master.Audio(video)
         fmt.Printf("%+v\n", audio)
      }
      fmt.Println(video)
   }
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
   file, err := os.Open("m3u8/cbc-segment.m3u8")
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

func TestBandwidth(t *testing.T) {
   master := &Master{Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 144},
      {Bandwidth: 1080},
      {Bandwidth: 720},
      {Bandwidth: 2160},
   }}
   sort.Sort(Bandwidth{master, 720})
   for _, str := range master.Stream {
      fmt.Println(str)
   }
}
