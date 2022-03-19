package main

import (
   "encoding/json"
   "flag"
   "os"
)

func doJSON(input, output string) error {
   src, err := os.Open(input)
   if err != nil {
      return err
   }
   defer src.Close()
   dst, err := os.Create(output)
   if err != nil {
      dst = os.Stdout
   }
   defer dst.Close()
   var value any
   if err := json.NewDecoder(src).Decode(&value); err != nil {
      return err
   }
   enc := json.NewEncoder(dst)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc.Encode(value)
}

func main() {
   // f
   var input string
   flag.StringVar(&input, "f", "", "input file")
   // o
   var output string
   flag.StringVar(&output, "o", "", "output file")
   flag.Parse()
   if input != "" {
      err := doJSON(input, output)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
