package main

import (
   "errors"
   "net/http"
   "os"
   "path"
)

const (
   zero = "https://vod-ap3-aoc.tv.apple.com/itunes-assets/HLSVideo126/v4/25/a3/dd/25a3ddc3-9b9e-cfc2-97e2-5dcb03ffd255/P377684155_A1524726231_FF_video_gr203_sdr_508x254_cbcs_--0.mp4"
   one = "https://vod-ap3-aoc.tv.apple.com/itunes-assets/HLSVideo126/v4/25/a3/dd/25a3ddc3-9b9e-cfc2-97e2-5dcb03ffd255/P377684155_A1524726231_FF_video_gr203_sdr_508x254_cbcs_--1.m4s"
)

func get(addr string) error {
   file, err := os.Create(path.Base(addr))
   if err != nil {
      return err
   }
   defer file.Close()
   res, err := http.Get(addr)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   if _, err := file.ReadFrom(res.Body); err != nil {
      return err
   }
   return nil
}

func main() {
   err := get(one)
   if err != nil {
      panic(err)
   }
}
