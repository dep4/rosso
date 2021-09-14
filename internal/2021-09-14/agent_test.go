package agent

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "net/http"
   "net/url"
   "os"
   "strings"
   "testing"
)

func TestAgent(t *testing.T) {
   pass := os.Getenv("ENCRYPTEDPASS")
   if pass == "" {
      t.Fatal("missing pass")
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
      t.Fatal(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   ua, err := os.Open("getAllUasJson.json")
   if err != nil {
      t.Fatal(err)
   }
   defer ua.Close()
   hash, err := os.Open("getAllHashesJson.json")
   if err != nil {
      t.Fatal(err)
   }
   defer hash.Close()
   j, err := ja3.NewJA3er(ua, hash)
   if err != nil {
      t.Fatal(err)
   }
   j.SortUsers()
   u, err := find(j, req)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(u)
}
