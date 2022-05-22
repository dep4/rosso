package format

import (
   "fmt"
   "io"
   "net/http"
   "os"
   "testing"
)

func TestString(t *testing.T) {
   tests := [][]byte{
      []byte("hello world 😀"),
      []byte("\xe0<\x00"),
      []byte{0, 1},
      []byte{0xE0, '<'},
   }
   for _, test := range tests {
      ok := IsString(test)
      fmt.Println(ok)
   }
}

func TestOpen(t *testing.T) {
   type token struct {
      Services string
      Token string
   }
   cache, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   tok, err := Open[token](cache, "googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", tok)
}

func TestLabel(t *testing.T) {
   fmt.Println(LabelNumber(9_999))
   fmt.Println(LabelSize(9_999))
   fmt.Println(LabelRate(9_999))
}

func TestProgress(t *testing.T) {
   res, err := http.Get("https://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := ProgressBytes(io.Discard, res.ContentLength)
   if _, err := io.Copy(pro, res.Body); err != nil {
      t.Fatal(err)
   }
}
