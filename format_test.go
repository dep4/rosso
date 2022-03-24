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

func TestProgress(t *testing.T) {
   res, err := http.Get("http://speedtest.lax.hivelocity.net/10Mio.dat")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := NewProgress(res)
   io.ReadAll(pro)
}

func TestLabel(t *testing.T) {
   fmt.Println(LabelNumber(9_999))
   fmt.Println(LabelSize(9_999))
   fmt.Println(LabelRate(9_999))
}

func TestPercent(t *testing.T) {
   fmt.Println(Percent(2, 3))
}
