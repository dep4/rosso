package main

import (
   "flag"
   "fmt"
   "github.com/89z/parse/http"
   "net/http/httputil"
   "os"
   stdhttp "net/http"
)

func main() {
   var (
      https bool
      output string
   )
   flag.BoolVar(&https, "s", false, "HTTPS")
   flag.StringVar(&output, "o", "", "output file")
   flag.Parse()
   if flag.NArg() != 1 {
      fmt.Println("http [flags] [request file]")
      flag.PrintDefaults()
      return
   }
   file := flag.Arg(0)
   rd, err := os.Open(file)
   if err != nil {
      panic(err)
   }
   defer rd.Close()
   req, err := http.ReadRequest(rd)
   if err != nil {
      panic(err)
   }
   if https {
      req.URL.Scheme = "https"
   } else {
      req.URL.Scheme = "http"
   }
   res, err := new(stdhttp.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   d, err := httputil.DumpResponse(res, false)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(d)
   if output == "" {
      os.Stdout.ReadFrom(res.Body)
      return
   }
   wr, err := os.Create(output)
   if err != nil {
      panic(err)
   }
   defer wr.Close()
   wr.ReadFrom(res.Body)
}
