package main

import (
   "fmt"
   "github.com/89z/format/hls"
   "net/http"
   "os"
   "sort"
)

func doManifest(address, output string, bandwidth int, info bool) error {
   fmt.Println("GET", address)
   res, err := http.Get(address)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   mas, err := hls.NewMaster(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   sort.Sort(hls.Bandwidth{mas, bandwidth})
   if info {
      for _, str := range mas.Stream {
         fmt.Println(str)
      }
   } else {
      for _, str := range mas.Stream {
         fmt.Println("GET", str.URI)
         res, err := http.Get(str.URI.String())
         if err != nil {
            return err
         }
         seg, err := hls.NewSegment(res.Request.URL, res.Body)
         if err != nil {
            return err
         }
         file, err := os.Create(output + seg.Ext())
         if err != nil {
            return err
         }
         for i, info := range seg.Info {
            fmt.Print(seg.Progress(i))
            res, err := http.Get(info.URI.String())
            if err != nil {
               return err
            }
            if _, err := file.ReadFrom(res.Body); err != nil {
               return err
            }
            if err := res.Body.Close(); err != nil {
               return err
            }
         }
         if err := res.Body.Close(); err != nil {
            return err
         }
         return file.Close()
      }
   }
   return nil
}
