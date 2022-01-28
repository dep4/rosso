package main

import (
   "flag"
   "fmt"
   "github.com/89z/format/net"
   "os"
)

func main() {
   var (
      https, info, redirect bool
      output string
   )
   flag.BoolVar(&info, "i", false, "info")
   flag.StringVar(&output, "o", "", "output file")
   flag.BoolVar(&redirect, "r", false, "redirect")
   flag.BoolVar(&https, "s", false, "HTTPS")
   flag.Parse()
   if flag.NArg() == 1 {
      input := flag.Arg(0)
      read, err := os.Open(input)
      if err != nil {
         panic(err)
      }
      defer read.Close()
      req, err := net.ReadRequest(read, https)
      if err != nil {
         panic(err)
      }
      file, err := os.Create(output)
      if err != nil {
         file = os.Stdout
      }
      defer file.Close()
      if info {
         err := net.WriteRequest(req, file)
         if err != nil {
            panic(err)
         }
      } else {
         res, err := roundTrip(req, redirect)
         if err != nil {
            panic(err)
         }
         defer res.Body.Close()
         if err := write(res, file); err != nil {
            panic(err)
         }
      }
   } else {
      fmt.Println("net [flags] [request file]")
      flag.PrintDefaults()
   }
}
