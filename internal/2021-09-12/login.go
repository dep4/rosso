package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "net/http"
)

var tests = []string{
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

func sanityCheck(hello string) error {
   for _, test := range tests {
      spec, err := ja3.Parse(hello)
      if err != nil {
         return err
      }
      req, err := http.NewRequest("HEAD", test, nil)
      if err != nil {
         return err
      }
      if _, err := ja3.NewTransport(spec).RoundTrip(req); err == nil {
         return nil
      }
   }
   return fmt.Errorf("manual review %q", hello)
}

func main() {
}
