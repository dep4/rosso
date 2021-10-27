package tls

import (
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

func TestHello(t *testing.T) {
   data, err := hex.DecodeString(android)
   if err != nil {
      t.Fatal(err)
   }
   hello, err := NewClientHello(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", hello.ClientHelloSpec)
   fmt.Printf("%#v\n", hello.Version)
}

func TestPcap(t *testing.T) {
   data, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      t.Fatal(err)
   }
   for _, hand := range Handshakes(data) {
      hello, err := NewClientHello(hand)
      if err == nil {
         fmt.Printf("%+v\n", hello)
      }
   }
}
