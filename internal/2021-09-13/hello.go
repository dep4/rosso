package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "net/http"
)

var tests = []string{
   "https://github.com",
   "https://stackoverflow.com",
   "https://variety.com",
   "https://vimeo.com",
   "https://www.google.com",
   "https://www.indiewire.com",
   "https://www.nytimes.com",
   "https://www.reddit.com",
   "https://www.wikipedia.org",
   "https://www.youtube.com",
}

const hello =
   "771," +
   "4866-4867-4865-49196-49200-49195-49199-52393-52392-159-158-52394-49327-49325-49326-49324-49188-49192-49187-49191-49162-49172-49161-49171-49315-49311-49314-49310-107-103-57-51-157-156-49313-49309-49312-49308-61-60-53-47-255," +
   "0-11-10-16-22-23-49-13-43-45-51-21,29-23-1035-25-24,0-1-2"

func main() {
   for _, test := range tests {
      spec, err := ja3.Parse(hello)
      if err != nil {
         panic(err)
      }
      req, err := http.NewRequest("HEAD", test, nil)
      if err != nil {
         panic(err)
      }
      _, err = ja3.NewTransport(spec).RoundTrip(req)
      fmt.Println(err, test)
   }
}
