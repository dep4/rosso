package hls

import (
   "net/http"
   "testing"
)

const addr = "http://s4b3b9a4.ssl.hwcdn.net/files/a8wn4hw/vi/04/07/10427421/hls-mi/s104274210.ts"

func Test_Ext(t *testing.T) {
   res, err := http.Get(addr)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   if err := read_file(res); err != nil {
      t.Fatal(err)
   }
}
