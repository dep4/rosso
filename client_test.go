package format

import (
   "fmt"
   "net/http"
   "testing"
)

func Test_Client(t *testing.T) {
   c := New_Client()
   req, err := http.NewRequest("HEAD", "http://godocs.io", nil)
   if err != nil {
      t.Fatal(err)
   }
   if _, err := c.Do(req); err != nil {
      fmt.Println(err)
   } else {
      t.Fatal(c)
   }
}
