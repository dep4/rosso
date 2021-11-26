package main

import (
   "bytes"
   "flag"
   "fmt"
   "github.com/89z/parse/net"
   "github.com/89z/parse/protobuf"
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
      https, proto, noBody bool
      output string
   )
   flag.BoolVar(&https, "s", false, "HTTPS")
   flag.BoolVar(&noBody, "n", false, "no body")
   flag.BoolVar(&proto, "p", false, "Protocol Buffer")
   flag.StringVar(&output, "o", "", "output file")
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
   if noBody {
      return
   }
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
      buf, err := mes.MarshalJSON()
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
