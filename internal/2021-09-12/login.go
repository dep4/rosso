package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "net/http"
   "net/url"
   "os"
   "strings"
   "time"
)

func main() {
   pass := os.Getenv("ENCRYPTEDPASS")
   if pass == "" {
      fmt.Println("missing pass")
      return
   }
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "sdk_version": {"17"},
      "EncryptedPasswd": {pass},
   }
   ua, err := os.Open("getAllUasJson.json")
   if err != nil {
      panic(err)
   }
   defer ua.Close()
   hash, err := os.Open("getAllHashesJson.json")
   if err != nil {
      panic(err)
   }
   defer hash.Close()
   j, err := ja3.NewJA3er(ua, hash)
   if err != nil {
      panic(err)
   }
   j.SortUsers()
   for _, user := range j.Users {
      fmt.Println(user)
      hello := j.JA3(user.MD5)
      spec, err := ja3.Parse(hello)
      if err != nil {
         panic(err)
      }
      req, err := http.NewRequest(
         "POST", "https://android.clients.google.com/auth",
         strings.NewReader(val.Encode()),
      )
      if err != nil {
         panic(err)
      }
      req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
      res, err := ja3.NewTransport(spec).RoundTrip(req)
      if err != nil {
         err := sanityCheck(hello)
         if err != nil {
            fmt.Println(err)
            continue
         }
      }
      defer res.Body.Close()
      fmt.Println(res.Status)
      if res.StatusCode == http.StatusOK {
         break
      }
      time.Sleep(time.Second)
   }
}

func sanityCheck(hello string) error {
   tests := []string{
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
