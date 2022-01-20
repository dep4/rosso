package format

import (
   "fmt"
   "io"
   "net/http"
   "testing"
)

func TestProgress(t *testing.T) {
   res, err := http.Get("http://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := NewProgress(res)
   io.ReadAll(pro)
}

func TestPercent(t *testing.T) {
   per := Percent(2, 3)
   fmt.Println(per)
}

func TestSymbol(t *testing.T) {
   nums := []int64{999, 1_234_567_890}
   for _, num := range nums {
      get := Number.GetInt64(num)
      fmt.Println(get)
   }
}
