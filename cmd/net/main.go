package main

import (
   "flag"
   "fmt"
   "github.com/89z/format/net"
   "os"
)

func main() {
   // f
   var name string
   flag.StringVar(&name, "f", "", "file")
   // i
   var info bool
   flag.BoolVar(&info, "i", false, "info")
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
      src, err := os.Open(name)
      if err != nil {
         panic(err)
      }
      defer src.Close()
      req, err := net.ReadRequest(src, https)
      if err != nil {
         panic(err)
      }
      dst, err := os.Create(output)
      if err != nil {
         dst = os.Stdout
      }
      defer dst.Close()
      if info {
         err := net.WriteRequest(req, dst)
         if err != nil {
            panic(err)
         }
      } else {
         res, err := roundTrip(req, redirect)
         if err != nil {
            panic(err)
         }
         defer res.Body.Close()
         if err := write(res, dst); err != nil {
            panic(err)
         }
      }
   } else {
      fmt.Println("net [flags]")
      flag.PrintDefaults()
   }
}
