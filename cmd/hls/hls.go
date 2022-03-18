package main

import (
   "flag"
   "fmt"
   "github.com/89z/format/hls"
   "github.com/89z/mech/nbc"
   "net/http"
   "os"
   "sort"
)

func doManifest(guid int64, bandwidth int, info bool) error {
   vod, err := nbc.NewAccessVOD(guid)
   if err != nil {
      return err
   }
   fmt.Println("GET", vod.ManifestPath)
   res, err := http.Get(vod.ManifestPath)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   mas, err := hls.NewMaster(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   if info {
      for _, str := range mas.Stream {
         fmt.Println(str)
      }
   } else {
      video, err := nbc.NewVideo(guid)
      if err != nil {
         return err
      }
      sort.Sort(hls.Bandwidth{mas, bandwidth})
      if err := download(video, mas.Stream[0]); err != nil {
         return err
      }
   }
   return nil
}

func download(video *nbc.Video, str hls.Stream) error {
   fmt.Println("GET", str.URI)
   res, err := http.Get(str.URI.String())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   seg, err := hls.NewSegment(res.Request.URL, res.Body)
   if err != nil {
      return err
   }
   file, err := os.Create(video.Base() + seg.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   fmt.Println(str)
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
   return nil
}

func main() {
   // b
   var guid int64
   flag.Int64Var(&guid, "b", 0, "GUID")
   // f
   var bandwidth int
   // nbc.com/saturday-night-live/video/march-12-zoe-kravitz/9000199371
   flag.IntVar(&bandwidth, "f", 5_581_000, "target bandwidth")
   // i
   var info bool
   flag.BoolVar(&info, "i", false, "info")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if verbose {
      nbc.LogLevel = 1
   }
   if guid >= 1 {
      err := doManifest(guid, bandwidth, info)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
