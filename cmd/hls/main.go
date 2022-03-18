package main

import (
   "flag"
)

func main() {
   // a
   var address string
   flag.StringVar(&address, "a", "", "address")
   // f
   var bandwidth int
   flag.IntVar(&bandwidth, "f", 0, "target bandwidth")
   // i
   var info bool
   flag.BoolVar(&info, "i", false, "info")
   // o
   var output string
   flag.StringVar(&output, "o", "output", "output")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if address != "" {
      err := doManifest(address, output, bandwidth, info)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
