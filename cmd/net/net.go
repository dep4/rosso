package main

import (
   "bytes"
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

func encode(r io.Reader, w io.Writer, iJSON, iProto bool) error {
   switch {
   case iJSON:
      src, err := io.ReadAll(r)
      if err != nil {
         return err
      }
      dst := new(bytes.Buffer)
      if err := json.Indent(dst, src, "", " "); err != nil {
         return err
      }
      if _, err := io.Copy(w, dst); err != nil {
         return err
      }
      return nil
   case iProto:
      src, err := io.ReadAll(r)
      if err != nil {
         return err
      }
      recs := protobuf.Bytes(src)
      enc := json.NewEncoder(w)
      enc.SetIndent("", " ")
      return enc.Encode(recs)
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
      https, iJSON, iProto, noBody bool
      output string
   )
   flag.BoolVar(&https, "s", false, "HTTPS")
   flag.BoolVar(&iJSON, "j", false, "indent JSON")
   flag.BoolVar(&iProto, "p", false, "indent Protocol Buffer")
   flag.BoolVar(&noBody, "n", false, "no body")
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
   if err := encode(res.Body, write, iJSON, iProto); err != nil {
      panic(err)
   }
}
