package main

import (
   "net/http"
   "os"
)

var addrs = []string{
   "https://vod-ap3-aoc.tv.apple.com/itunes-assets/HLSVideo116/v4/20/b0/75/20b075e2-f120-1929-7362-edbf5715b8e0/P377684155_A1524726231_audio_en_gr160-0.mp4",
   "https://vod-ap3-aoc.tv.apple.com/itunes-assets/HLSVideo116/v4/20/b0/75/20b075e2-f120-1929-7362-edbf5715b8e0/P377684155_A1524726231_audio_en_gr160-1.m4s",
}

func main() {
   file, err := os.Create("enc.mp4")
   if err != nil {
      panic(err)
   }
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
