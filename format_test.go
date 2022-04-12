package format

import (
   "fmt"
   "io"
   "net/http"
   "testing"
)

type token struct {
   Services string
   Token string
}

func TestDecode(t *testing.T) {
   tok, err := Open[token]("token.json")
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

func TestPercent(t *testing.T) {
   fmt.Println(Percent(2, 3))
}

func TestProgress(t *testing.T) {
   pro := NewProgress(io.Discard, 1)
   res, err := http.Get("https://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   if _, err := pro.Copy(res); err != nil {
      t.Fatal(err)
   }
}
