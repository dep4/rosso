package main

import (
   "flag"
   "github.com/89z/std/http"
)

var client = http.Default_Client

type flags struct {
   address string
   bandwidth_audio int
   bandwidth_video int
   info bool
   key string
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.IntVar(&f.bandwidth_video, "f", 1_999_999, "video bandwidth")
   flag.IntVar(&f.bandwidth_audio, "g", 127_000, "audio bandwidth")
   flag.BoolVar(&f.info, "i", false, "information")
   flag.StringVar(&f.key, "k", "", "key")
   flag.Parse()
   if f.address != "" {
      err := f.DASH()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
