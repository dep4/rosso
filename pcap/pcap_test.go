package pcap

import (
   "fmt"
   "os"
   "testing"
)

func TestPcap(t *testing.T) {
   data, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      t.Fatal(err)
   }
   for _, hand := range Handshakes(data) {
      spec, err := hand.ClientHello()
      if err == nil {
         fmt.Println(hand)
         fmt.Printf("%+v\n", spec)
      }
   }
}
