package main

import (
   "flag"
   "github.com/89z/std/dash"
   "github.com/89z/std/http"
)

var client = http.Default_Client

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.IntVar(&f.bandwidth_video, "f", 1, "video bandwidth")
   flag.IntVar(&f.bandwidth_audio, "g", 1, "audio bandwidth")
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

type flags struct {
   address string
   bandwidth_audio int
   bandwidth_video int
   info bool
   key string
}

type stream struct {
   bandwidth int
   base string
   dash.Representations
   key []byte
}
