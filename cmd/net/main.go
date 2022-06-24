package main

import (
   "flag"
   "github.com/89z/format"
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
   // r
   var redirect bool
   flag.BoolVar(&redirect, "r", false, "redirect")
   // s
   var https bool
   flag.BoolVar(&https, "s", false, "HTTPS")
   flag.Parse()
   if name != "" {
      out, err := format.Create(output)
      if err != nil {
         out = os.Stdout
      }
      defer out.Close()
      src, err := os.Open(name)
      if err != nil {
         panic(err)
      }
      defer src.Close()
      req, err := net.Read_Request(src)
      if err != nil {
         panic(err)
      }
      if req.URL.Scheme == "" {
         if https {
            req.URL.Scheme = "https"
         } else {
            req.URL.Scheme = "http"
         }
      }
      if golang {
         err := net.Write_Request(req, out)
         if err != nil {
            panic(err)
         }
      } else {
         err := write(req, redirect, out)
         if err != nil {
            panic(err)
         }
      }
   } else {
      flag.Usage()
   }
}
