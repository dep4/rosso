package hls

import (
   "encoding/json"
   "fmt"
   "net/http"
   "os"
   "path"
   "strings"
   "testing"
)

func TestSegment(t *testing.T) {
   mas, err := newMaster()
   if err != nil {
      t.Fatal(err)
   }
   str := mas.GetStream(func(s Stream) bool {
      return s.Bandwidth < 400_000
   })
   uris := []string{str.URI, mas.GetMedia(str).URI}
   for _, uri := range uris {
      res, err := http.Get(uri)
      if err != nil {
         t.Fatal(err)
      }
      defer res.Body.Close()
      seg, err := NewSegment(res.Request.URL, res.Body)
      if err != nil {
         t.Fatal(err)
      }
      if err := decrypt(seg); err != nil {
         t.Fatal(err)
      }
   }
}

func decrypt(seg *Segment) error {
   dec, err := doKey(seg)
   if err != nil {
      return err
   }
   res, err := http.Get(seg.Info[0].URI)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   file, err := os.Create("ignore/" + path.Base(seg.Info[0].URI))
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := dec.Copy(file, res.Body); err != nil {
      return err
   }
   return nil
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


func doKey(seg *Segment) (*Decrypter, error) {
   res, err := http.Get(seg.Key.URI)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return NewDecrypter(res.Body)
}


