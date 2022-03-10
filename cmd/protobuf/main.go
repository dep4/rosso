package main

import (
   "flag"
)

func main() {
   // f
   var input string
   flag.StringVar(&input, "f", "", "input file")
   // o
   var output string
   flag.StringVar(&output, "o", "", "output file")
   flag.Parse()
   if input != "" {
      err := doProtoBuf(input, output)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
