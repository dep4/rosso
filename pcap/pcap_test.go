package pcap

import (
   "fmt"
   "os"
   "testing"
)

func TestPcap(t *testing.T) {
   f, err := os.Open("PCAPdroid_22_Oct_15_19_28.pcap")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   hands, err := Handshakes(f)
   if err != nil {
      t.Fatal(err)
   }
   for _, hand := range hands {
      fmt.Println(hand)
   }
}
