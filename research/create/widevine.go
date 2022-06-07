package main

import (
   "os"
   "net/http"
)

var addrs = []string{
   "https://vod-ap3-aoc.tv.apple.com/itunes-assets/HLSVideo126/v4/25/a3/dd/25a3ddc3-9b9e-cfc2-97e2-5dcb03ffd255/P377684155_A1524726231_FF_video_gr203_sdr_508x254_cbcs_--0.mp4",
   "https://vod-ap3-aoc.tv.apple.com/itunes-assets/HLSVideo126/v4/25/a3/dd/25a3ddc3-9b9e-cfc2-97e2-5dcb03ffd255/P377684155_A1524726231_FF_video_gr203_sdr_508x254_cbcs_--1.m4s",
}

func main() {
   file, err := os.Create("ignore/enc.mp4")
   if err != nil {
      panic(err)
   }
   defer file.Close()
   for _, addr := range addrs {
      res, err := http.Get(addr)
      if err != nil {
         panic(err)
      }
      if res.StatusCode != http.StatusOK {
         panic(res.Status)
      }
      if _, err := file.ReadFrom(res.Body); err != nil {
         panic(err)
      }
      if err := res.Body.Close(); err != nil {
         panic(err)
      }
   }
}
