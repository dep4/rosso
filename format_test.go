package format

import (
   "fmt"
   "io"
   "net/http"
   "testing"
)

func Test_String(t *testing.T) {
   tests := [][]byte{
      []byte("hello world 😀"),
      []byte("\xe0<\x00"),
      []byte{0, 1},
      []byte{0xE0, '<'},
   }
   for _, test := range tests {
      ok := String(test)
      fmt.Println(ok)
   }
}

func Test_Progress(t *testing.T) {
   res, err := http.Get("https://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := Progress_Bytes(io.Discard, res.ContentLength)
   if _, err := io.Copy(pro, res.Body); err != nil {
      t.Fatal(err)
   }
}
