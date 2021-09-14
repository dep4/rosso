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

func callback(hello string) error {
   pass := os.Getenv("ENCRYPTEDPASS")
   if pass == "" {
      return fmt.Errorf("ENCRYPTEDPASS %q", pass)
   }
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "sdk_version": {"17"},
      "EncryptedPasswd": {pass},
   }
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth",
      strings.NewReader(val.Encode()),
   )
   if err != nil {
      return err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   spec, err := ja3.Parse(hello)
   if err != nil {
      return err
   }
   time.Sleep(100 * time.Millisecond)
   res, err := ja3.NewTransport(spec).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return fmt.Errorf("status %q", res.Status)
   }
   return nil
}

func main() {
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
   user := j.Find(callback)
   fmt.Println(user)
}
