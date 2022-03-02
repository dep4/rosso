package hls

import (
   "fmt"
   "net/http"
   "os"
   "testing"
)

func doSegment() (*Segment, error) {
   mas, err := doMaster()
   if err != nil {
      return nil, err
   }
   res, err := http.Get(mas.Stream[0].URI)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return NewSegment(res.Request.URL, res.Body)
}

func doKey(seg *Segment) (*Decrypter, error) {
   res, err := http.Get(seg.Key.URI)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return NewDecrypter(res.Body)
}

func TestSegment(t *testing.T) {
   seg, err := doSegment()
   if err != nil {
      t.Fatal(err)
   }
   dec, err := doKey(seg)
   if err != nil {
      t.Fatal(err)
   }
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
   res, err := http.Get(seg.Info[1].URI)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   buf, err := dec.Decrypt(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   if err := os.WriteFile("ignore.ts", buf, os.ModePerm); err != nil {
      t.Fatal(err)
   }
}


