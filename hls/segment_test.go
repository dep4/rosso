package hls

import (
   "fmt"
   "net/http"
   "testing"
)

const pbs = "https://urs.pbs.org/redirect/2dc8ce48e5d54172ad141e078d04cc4d/"

func TestSegment(t *testing.T) {
   master, err := newMaster()
   if err != nil {
      t.Fatal(err)
   }
   addr := master.Stream[0].URI
   res, err := http.Get(addr.String())
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   seg, err := NewSegment(res.Request.URL, res.Body)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", seg)
}

func newMaster() (*Master, error) {
   res, err := http.Get(pbs)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return NewMaster(res.Request.URL, res.Body)
}
