package main

import (
   "flag"
   "github.com/89z/format/net"
   "os"
)

func main() {
   // f
   var name string
   flag.StringVar(&name, "f", "", "input file")
   // g
   var golang bool
   flag.BoolVar(&golang, "g", false, "request as Go code")
   // o
   var output string
   flag.StringVar(&output, "o", "", "output file")
   // s
   var scheme string
   flag.StringVar(&scheme, "s", "http", "scheme")
   flag.Parse()
   if name != "" {
      dst, err := os.Create(output)
      if err != nil {
         dst = os.Stdout
      }
      defer dst.Close()
      src, err := os.Open(name)
      if err != nil {
         panic(err)
      }
      defer src.Close()
      req, err := net.ReadRequest(src)
      if err != nil {
         panic(err)
      }
      if req.URL.Scheme == "" {
         req.URL.Scheme = scheme
      }
      if golang {
         err := net.WriteRequest(req, dst)
         if err != nil {
            panic(err)
         }
      } else {
         err := write(req, dst)
         if err != nil {
            panic(err)
         }
      }
   } else {
      flag.Usage()
   }
}
