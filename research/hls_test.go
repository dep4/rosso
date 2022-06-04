package hls

import (
   "fmt"
   "net/url"
   "os"
   "sort"
   "testing"
)

func TestMaster(t *testing.T) {
   file, err := os.Open("ignore.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := NewScanner(file).Master(&url.URL{})
   if err != nil {
      t.Fatal(err)
   }
   var streams []Stream
   for stream := range master.Streams {
      streams = append(streams, stream)
   }
   sort.Slice(streams, func(a, b int) bool {
      return streams[a].Bandwidth < streams[b].Bandwidth
   })
   for _, stream := range streams {
      fmt.Println(stream)
   }
}
