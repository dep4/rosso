package hls

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "sort"
   "testing"
)

const pbs = "https://urs.pbs.org/redirect/2dc8ce48e5d54172ad141e078d04cc4d/"

func TestMaster(t *testing.T) {
   res, err := http.Get(pbs)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   master, err := NewMaster(res.Request.URL, res.Body)
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Stream {
      fmt.Println(stream)
   }
   media := master.GetMedia(master.Stream[0])
   fmt.Printf("%+v\n", media)
}

func TestProgress(t *testing.T) {
   seg := Segment{
      Info: make([]Information, 9),
   }
   for i := range seg.Info {
      fmt.Print(seg.Progress(i))
   }
   fmt.Println("END")
}

func TestSort(t *testing.T) {
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

var ivs = []string{
   "00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, iv := range ivs {
      buf, err := hexDecode(iv)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(buf)
   }
}

func TestSegment(t *testing.T) {
   seg, err := newSegment()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(seg.Key.URI.String())
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   block, err := NewCipher(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.ts")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
      res, err := http.Get(info.URI.String())
      if err != nil {
         t.Fatal(err)
      }
      defer res.Body.Close()
      buf, err := block.Decrypt(info, res.Body)
      if err != nil {
         t.Fatal(err)
      }
      if _, err := file.Write(buf); err != nil {
         t.Fatal(err)
      }
   }
}

func newSegment() (*Segment, error) {
   file, err := os.Open("m3u8/abc-segment.m3u8")
   if err != nil {
      return nil, err
   }
   defer file.Close()
   return NewSegment(&url.URL{}, file)
}
