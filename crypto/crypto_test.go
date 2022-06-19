package crypto

import (
   "fmt"
   "net/http"
   "testing"
)

func Test_Format_JA3(t *testing.T) {
   hello, err := Parse_JA3(Android_API_26)
   if err != nil {
      t.Fatal(err)
   }
   ja3, err := Format_JA3(hello)
   if err != nil {
      t.Fatal(err)
   }
   if ja3 != Android_API_26 {
      t.Fatal(ja3)
   }
}

func Test_Transport(t *testing.T) {
   req, err := http.NewRequest("HEAD", "https://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   hello, err := Parse_JA3(Android_API_26)
   if err != nil {
      t.Fatal(err)
   }
   res, err := Transport(hello).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", res)
}
