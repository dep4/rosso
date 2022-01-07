package format

import (
   "fmt"
   "io"
   "net/http"
   "os"
   "testing"
)

func TestProgress(t *testing.T) {
   res, err := http.Get("http://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := NewProgress(res, os.Stdout)
   io.ReadAll(pro)
}

func TestPercent(t *testing.T) {
   Percent(os.Stdout, 2, 3)
   fmt.Println()
}

func TestSymbol(t *testing.T) {
   nums := []int64{999, 1_234_567_890}
   for _, num := range nums {
      Number.Int64(os.Stdout, num)
      fmt.Println()
   }
}
