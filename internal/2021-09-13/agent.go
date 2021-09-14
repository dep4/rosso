package main

import (
   "bytes"
   "fmt"
   "github.com/89z/parse/ja3"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
   "time"
)

func body(r io.ReadCloser) ([]byte, error) {
   if r == nil {
      return nil, nil
   }
   return io.ReadAll(r)
}

func a(j *ja3.JA3er, req *http.Request) (string, error) {
   j.SortUsers()
   data, err := body(req.Body)
   if err != nil {
      return "", err
   }
   done := make(map[string]bool)
   for _, user := range j.Users {
      hello := j.JA3(user.MD5)
      if done[hello] {
         continue
      } else {
         done[hello] = true
      }
      spec, err := ja3.Parse(hello)
      if err != nil {
         fmt.Println(err)
         continue
      }
      req.Body = io.NopCloser(bytes.NewReader(data))
      time.Sleep(100 * time.Millisecond)
      res, err := ja3.NewTransport(spec).RoundTrip(req)
      if err != nil {
         fmt.Println(err)
         continue
      }
      defer res.Body.Close()
      fmt.Println(res.Status)
      if res.StatusCode == http.StatusOK {
         return hello, nil
      }
   }
   return "", fmt.Errorf("%+v FAIL", req)
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
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth",
      strings.NewReader(val.Encode()),
   )
   if err != nil {
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   s, err := a(j, req)
   if err != nil {
      panic(err)
   }
   fmt.Println(s)
}
