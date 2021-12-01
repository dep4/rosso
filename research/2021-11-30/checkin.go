package main

import (
   "bytes"
   "fmt"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
)

var mes = protobuf.Message{
   3:"1-da39a3ee5e6b4b0d3255bfef95601890afd80709",
   4:protobuf.Message{},
}

var androidKey = []byte("AAAAgMom")

func main() {
   req, err := http.NewRequest(
      "POST", "http://android.clients.google.com/checkin", mes.Encode(),
   )
   if err != nil {
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      panic(err)
   }
   if bytes.Contains(buf, androidKey) {
      fmt.Println("pass")
   } else {
      fmt.Println("fail")
   }
}
