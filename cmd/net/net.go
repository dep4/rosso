package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "github.com/89z/parse/net"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
   "net/http/httputil"
   "os"
)

func encode(r io.Reader, w io.Writer, proto bool) error {
   if proto {
      buf, err := io.ReadAll(r)
      if err != nil {
         return err
      }
      dec := protobuf.NewDecoder(buf)
      enc := json.NewEncoder(w)
      enc.SetIndent("", " ")
      return enc.Encode(dec)
   }
   _, err := io.Copy(w, r)
   if err != nil {
      return err
   }
   return nil
}

func file(output string) (*os.File, error) {
   if output != "" {
      return os.Create(output)
   }
   return os.Stdout, nil
}

func main() {
   var (
      https, proto bool
      output string
   )
   flag.BoolVar(&https, "s", false, "HTTPS")
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
   write, err := file(output)
   if err != nil {
      panic(err)
   }
   defer write.Close()
   if err := encode(res.Body, write, proto); err != nil {
      panic(err)
   }
}
