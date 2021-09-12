package main

import (
   "github.com/89z/parse/ja3"
   "net/http"
)

const test =
   "771," +
   "49196-49195-49200-49199-159-158-49188-49187-49192-49191-49162-49161-49172-49171-157-156-61-60-53-47-10," +
   "0-10-11-13-35-23-65281,,"

func main() {
   spec, err := ja3.Parse(test)
   if err != nil {
      panic(err)
   }
   req, err := http.NewRequest("HEAD", "https://example.com", nil)
   if err != nil {
      panic(err)
   }
   if _, err := ja3.NewTransport(spec).RoundTrip(req); err != nil {
      panic(err)
   }
}
