package main

import (
   "flag"
   "fmt"
   "github.com/89z/mech/html"
   "os"
)

func newFile(name string) (*os.File, error) {
   if name == "" {
      return os.Stdout, nil
   }
   return os.Create(name)
}

func main() {
   var indent, output string
   flag.StringVar(&indent, "i", "", "indent")
   flag.StringVar(&output, "o", "", "output file")
   flag.Parse()
   if flag.NArg() != 1 {
      fmt.Println("html-indent [flags] [input file]")
      flag.PrintDefaults()
      return
   }
   input := flag.Arg(0)
   r, err := os.Open(input)
   if err != nil {
      panic(err)
   }
   defer r.Close()
   w, err := newFile(output)
   if err != nil {
      panic(err)
   }
   defer w.Close()
   if err := html.NewLexer(r).Render(w, indent); err != nil {
      panic(err)
   }
}
