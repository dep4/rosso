package main

import (
   "github.com/89z/format"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

func round_trip(req *http.Request, redirect bool) (*http.Response, error) {
   if redirect {
      return new(http.Client).Do(req)
   }
   return new(http.Transport).RoundTrip(req)
}

func write(req *http.Request, redirect bool, file *os.File) error {
   res, err := round_trip(req, redirect)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if file == os.Stdout {
      buf, err := httputil.DumpResponse(res, true)
      if err != nil {
         return err
      }
      if format.Is_String(buf) {
         file.Write(buf)
      } else {
         quote := strconv.Quote(string(buf))
         file.WriteString(quote)
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
