package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "net/http"
   "time"
)

var sites = []string{
   "https://www.reddit.com",
   "https://github.com",
   "https://nebulance.io",
   "https://stackoverflow.com",
   "https://variety.com",
   "https://vimeo.com",
   "https://www.google.com",
   "https://www.indiewire.com",
   "https://www.wikipedia.org",
   "https://www.youtube.com",
}

func main() {
   for _, hello := range hellos {
      fmt.Println(hello + ":")
      for _, site := range sites {
         spec, err := ja3.Parse(hello)
         if err != nil {
            panic(err)
         }
         req, err := http.NewRequest("HEAD", site, nil)
         if err != nil {
            panic(err)
         }
         _, err = ja3.NewTransport(spec).RoundTrip(req)
         fmt.Println(err, site)
      }
      time.Sleep(time.Second)
   }
}

func version(min uint16) []uint16 {
   vs := []uint16{772, 771, 770, 769, 768}
   for k, v := range vs {
      if v == min {
         return vs[:k+1]
      }
   }
   return nil
}
