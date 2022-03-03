package hls

import (
   "encoding/json"
   "fmt"
   "net/http"
   "os"
   "strings"
   "testing"
)

func newTopaz() (string, error) {
   var buf strings.Builder
   buf.WriteString("https://topaz.viacomcbs.digital/topaz/api/")
   buf.WriteString("mgid:arc:showvideo:mtv.com:d26f2b22-097d-11e3-8a73-0026b9414f30/")
   buf.WriteString("mica.json?clientPlatform=android")
   res, err := http.Get(buf.String())
   if err != nil {
      return "", err
   }
   defer res.Body.Close()
   var topaz struct {
      StitchedStream struct {
         Source string
      }
   }
   if err := json.NewDecoder(res.Body).Decode(&topaz); err != nil {
      return "", err
   }
   return topaz.StitchedStream.Source, nil
}

func newMaster() (*Master, error) {
   topaz, err := newTopaz()
   if err != nil {
      return nil, err
   }
   res, err := http.Get(topaz)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return NewMaster(res.Request.URL, res.Body)
}

func TestMaster(t *testing.T) {
   mas, err := newMaster()
   if err != nil {
      t.Fatal(err)
   }
   for _, med := range mas.Media {
      fmt.Printf("%+v\n", med)
   }
   for _, str := range mas.Stream {
      fmt.Println(str)
   }
}

func doKey(seg *Segment) (*Decrypter, error) {
   res, err := http.Get(seg.Key.URI)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return NewDecrypter(res.Body)
}

func newSegment() (*Segment, error) {
   mas, err := newMaster()
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

func TestSegment(t *testing.T) {
   seg, err := newSegment()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", seg.Key)
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
   dec, err := doKey(seg)
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(seg.Info[0].URI)
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
