package main

import (
   "bytes"
   "flag"
   "encoding/json"
   "fmt"
   "github.com/89z/format/net"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/http/httputil"
   "os"
)

func file(output string) (*os.File, error) {
   if output != "" {
      return os.Create(output)
   }
   return os.Stdout, nil
}

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
   req, err := net.ReadRequest(read)
   if err != nil {
      panic(err)
   }
   if info {
      fmt.Printf("%#v\n", req.URL.Query())
      fmt.Printf("%#v\n", req.Header)
   }
   if https {
      req.URL.Scheme = "https"
   } else {
      req.URL.Scheme = "http"
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
   // body
   write, err := file(output)
   if err != nil {
      panic(err)
   }
   defer write.Close()
   if proto {
      mes, err := protobuf.Decode(res.Body)
      if err != nil {
         panic(err)
      }
      buf, err := json.Marshal(mes)
      if err != nil {
         panic(err)
      }
      if _, err := write.ReadFrom(bytes.NewReader(buf)); err != nil {
         panic(err)
      }
   } else {
      _, err := write.ReadFrom(res.Body)
      if err != nil {
         panic(err)
      }
   }
}
