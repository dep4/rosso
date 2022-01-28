package main

import (
   "github.com/89z/format"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

func roundTrip(req *http.Request, redirect bool) (*http.Response, error) {
   if redirect {
      return new(http.Client).Do(req)
   }
   return new(http.Transport).RoundTrip(req)
}

func write(res *http.Response, file *os.File) error {
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
