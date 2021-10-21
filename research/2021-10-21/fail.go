package main

import (
   "net/http"
   "net/http/httputil"
   "os"
)

func main() {
   req, err := http.NewRequest(
      "GET", "https://www.amazon.com/dp/B07K5214NZ", nil,
   )
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   dum, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(dum)
}
