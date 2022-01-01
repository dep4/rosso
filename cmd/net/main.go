package main

import (
   "flag"
   "fmt"
   "github.com/89z/format/net"
   "net/http"
   "net/http/httputil"
   "os"
)

func main() {
   var (
      https, info, proto bool
      output string
   )
   flag.BoolVar(&info, "i", false, "info")
   flag.StringVar(&output, "o", "", "output file")
   flag.BoolVar(&proto, "p", false, "Protocol Buffer")
   flag.BoolVar(&https, "s", false, "HTTPS")
   flag.Parse()
   if flag.NArg() != 1 {
      fmt.Println("net [flags] [request file]")
      flag.PrintDefaults()
      return
   }
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
   if info {
      fmt.Printf("%#v\n", req.URL.Query())
      fmt.Printf("%#v\n", req.Header)
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   // head
   buf, err := httputil.DumpResponse(res, false)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
   if err := readFrom(res.Body, output, proto); err != nil {
      panic(err)
   }
}
