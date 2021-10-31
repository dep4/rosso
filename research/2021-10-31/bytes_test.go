package bytes

import (
   "fmt"
   "os"
   "testing"
)

func TestHandshake(t *testing.T) {
   data, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      t.Fatal(err)
   }
   hand := handshake(data)
   fmt.Printf("%+v\n", hand)
}
