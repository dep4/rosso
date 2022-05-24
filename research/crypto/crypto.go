package main

import (
   "github.com/89z/format/crypto"
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
)

//pass
const ja3 = "771,49195-49196-52393-49199-49200-52392-49161-49162-49171-49172-156-157-47-53,65281-0-23-35-13-5-16-11-10,29-23-24,0"

//fail
//const ja3 = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513-21,29-23-24,0"

func main() {
   req := new(http.Request)
   req.Body = io.NopCloser(body)
   req.Header = make(http.Header)
   req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
   req.Method = "POST"
   req.URL = new(url.URL)
   req.URL.Host = "android.googleapis.com"
   req.URL.Path = "/auth"
   req.URL.Scheme = "https"
   hello, err := crypto.ParseJA3(ja3)
   if err != nil {
      panic(err)
   }
   res, err := crypto.Transport(hello).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
}
