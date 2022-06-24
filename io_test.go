package format

import (
   "fmt"
   "net/http"
   "testing"
)

func Test_Client(t *testing.T) {
   req, err := http.NewRequest("HEAD", "http://godocs.io", nil)
   if err != nil {
      t.Fatal(err)
   }
   if _, err := Default_Client.Do(req); err != nil {
      fmt.Println(err)
   } else {
      t.Fatal(req)
   }
}
