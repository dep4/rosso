package format

import (
   "fmt"
   "io"
   "net/http"
   "os"
   "testing"
)

type token struct {
   Services string
   Token string
}

func TestOpen(t *testing.T) {
   cache, err := os.UserCacheDir()
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
