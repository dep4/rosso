package hls

import (
   "io"
   "net/http"
   "strconv"
   "text/scanner"
   "unicode"
)

func one(req *http.Request) (*response, error) {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   return &response{res}, nil
}

type response struct {
   *http.Response
}
