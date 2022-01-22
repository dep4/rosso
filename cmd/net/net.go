package main

import (
   "flag"
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/net"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

func main() {
   var (
      https, info bool
      output string
   )
   flag.BoolVar(&info, "i", false, "info")
   flag.StringVar(&output, "o", "", "output file")
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
         err := roundTrip(req, file)
         if err != nil {
            panic(err)
         }
      }
   } else {
      fmt.Println("net [flags] [request file]")
      flag.PrintDefaults()
   }
}

func roundTrip(req *http.Request, file *os.File) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if file == os.Stdout {
      buf, err := httputil.DumpResponse(res, true)
      if err != nil {
         return err
      }
      if format.IsBinary(buf) {
         quote := strconv.Quote(string(buf))
         file.WriteString(quote)
      } else {
         file.Write(buf)
      }
   } else {
      buf, err := httputil.DumpResponse(res, false)
      if err != nil {
         return err
      }
      os.Stdout.Write(buf)
      if _, err := file.ReadFrom(res.Body); err != nil {
         return err
      }
   }
   return nil
}
