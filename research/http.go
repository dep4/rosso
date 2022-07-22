package http

import (
   "net/http"
)

func Clone(req *http.Request) *http.Request {
   req2 := new(http.Request)
   *req2 = *req
   return req2
}
