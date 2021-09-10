package ja3

import (
   "fmt"
   "net/http"
   "os"
   "testing"
)

func TestAgent(t *testing.T) {
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
   j, err := NewJA3er(ua, hash)
   if err != nil {
      t.Fatal(err)
   }
   j.SortUsers()
   user := j.Users[0]
   fmt.Printf("%+v\n", user)
   spec, err := Parse(j.JA3(user.MD5))
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest("GET", "https://ja3er.com/json", nil)
   if err != nil {
      t.Fatal(err)
   }
   res, err := NewTransport(spec).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
}
