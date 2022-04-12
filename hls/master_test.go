package hls

import (
   "fmt"
   "net/http"
   "testing"
)

const pbs = "https://urs.pbs.org/redirect/2dc8ce48e5d54172ad141e078d04cc4d/"

func TestMaster(t *testing.T) {
   res, err := http.Get(pbs)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   master, err := NewScanner(res.Body).Master(res.Request.URL)
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Stream {
      fmt.Println(stream)
   }
   media := master.GetMedia(master.Stream[0])
   fmt.Printf("%+v\n", media)
}
