package main

import (
   "github.com/89z/format"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

func write(req *http.Request, file *os.File) error {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
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
   return res.Body.Close()
}
