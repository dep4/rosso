package main

import (
   "encoding/json"
   "flag"
   "os"
)

func do_json(input, output string) error {
   in, err := os.Open(input)
   if err != nil {
      return err
   }
   defer in.Close()
   out, err := os.Create(output)
   if err != nil {
      out = os.Stdout
   }
   defer out.Close()
   var value any
   if err := json.NewDecoder(in).Decode(&value); err != nil {
      return err
   }
   enc := json.NewEncoder(out)
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
      err := do_json(input, output)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
