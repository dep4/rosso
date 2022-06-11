package format

import (
   "encoding/json"
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
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   file, err := Open(home, "googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   var token struct {
      Services string
      Token string
   }
   json.NewDecoder(file).Decode(&token)
   fmt.Printf("%+v\n", token)
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
