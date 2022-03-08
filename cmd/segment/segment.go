package main

import (
   "fmt"
   "github.com/89z/format/hls"
   "net/http"
   "os"
)

func main() {
   if len(os.Args) == 2 {
      addr := os.Args[1]
      fmt.Println("GET", addr)
      res, err := http.Get(addr)
      if err != nil {
         panic(err)
      }
      defer res.Body.Close()
      seg, err := hls.NewSegment(res.Request.URL, res.Body)
      if err != nil {
         panic(err)
      }
      if err := download(seg); err != nil {
         panic(err)
      }
   } else {
      fmt.Println("segment [URL]")
   }
}

func download(seg *hls.Segment) error {
   fmt.Println("GET", seg.Key.URI)
   res, err := http.Get(seg.Key.URI.String())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   dec, err := hls.NewDecrypter(res.Body)
   if err != nil {
      return err
   }
   file, err := os.Create("outfile" + seg.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   for i, info := range seg.Info {
      if i >= 1 {
         fmt.Print(" ")
      }
      fmt.Print(len(seg.Info)-i)
      res, err := http.Get(info.URI.String())
      if err != nil {
         return err
      }
      if _, err := dec.Copy(file, res.Body); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return nil
}
