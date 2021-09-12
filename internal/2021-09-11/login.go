package main

import (
   "fmt"
   "github.com/89z/ja3"
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
      spec, err := ja3.Parse(j.JA3(user.MD5))
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
         panic(err)
      }
      defer res.Body.Close()
      fmt.Println(res.Status)
      if res.StatusCode == http.StatusOK {
         break
      }
      time.Sleep(time.Second)
   }
}
