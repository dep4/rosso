package main

import (
   "flag"
   "fmt"
   "github.com/89z/format/hls"
   "net/http"
   "os"
   "sort"
)

func main() {
   // a
   var address string
   flag.StringVar(&address, "a", "", "address")
   // f
   var bandwidth int64
   flag.Int64Var(&bandwidth, "f", 420_000, "min bandwidth")
   // i
   var info bool
   flag.BoolVar(&info, "i", false, "info")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if verbose {
      mtv.LogLevel = 1
   }
   if address != "" {
      err := doManifest(address, bandwidth, info)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func doManifest(addr string, bandwidth int64, info bool) error {
   prop, err := mtv.NewItem(addr).Property()
   if err != nil {
      return err
   }
   top, err := prop.Topaz()
   if err != nil {
      return err
   }
   fmt.Println("GET", top.StitchedStream.Source)
   res, err := http.Get(top.StitchedStream.Source)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   mas, err := hls.NewMaster(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   sort.Slice(mas.Stream, func(a, b int) bool {
      return mas.Stream[a].Bandwidth < mas.Stream[b].Bandwidth
   })
   if info {
      for _, str := range mas.Stream {
         str.URI = nil
         fmt.Println(str)
      }
   } else {
      stream := mas.GetStream(func(s hls.Stream) bool {
         return s.Bandwidth >= bandwidth
      })
      err := download(stream, prop)
      if err != nil {
         return err
      }
   }
   return nil
}

func newSegment(str *hls.Stream) (*hls.Segment, error) {
   fmt.Println("GET", str.URI)
   res, err := http.Get(str.URI.String())
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return hls.NewSegment(res.Request.URL, res.Body)
}

func download(str *hls.Stream, prop *mtv.Property) error {
   seg, err := newSegment(str)
   if err != nil {
      return err
   }
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
   file, err := os.Create(prop.Base() + seg.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   for i, info := range seg.Info {
      fmt.Println(i, len(seg.Info)-1)
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
