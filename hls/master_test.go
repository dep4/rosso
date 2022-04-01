package hls

import (
   "fmt"
   "net/http"
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
