package tls

import (
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

const android = "16030100bb010000b703034420d198e7852decbc117dc7f90550b98f2d643c954bf3361ddaf127ff921b04000024c02bc02ccca9c02fc030cca8009e009fc009c00ac013c01400330039009c009d002f00350100006aff0100010000000022002000001d636c69656e7473657276696365732e676f6f676c65617069732e636f6d0017000000230000000d0016001406010603050105030401040303010303020102030010000b000908687474702f312e31000b00020100000a000400020017"

func TestHandshakes(t *testing.T) {
   data, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      t.Fatal(err)
   }
   for _, hand := range handshakes(data) {
      hello, err := ParseHandshake(hand)
      if err == nil {
         fmt.Printf("%+v\n", hello)
      }
   }
}

func TestHandshake(t *testing.T) {
   data, err := hex.DecodeString(android)
   if err != nil {
      t.Fatal(err)
   }
   hello, err := ParseHandshake(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", hello)
}

func TestJA3(t *testing.T) {
   h, err := ParseJA3(Android)
   if err != nil {
      t.Fatal(err)
   }
   j, err := h.FormatJA3()
   if err != nil {
      t.Fatal(err)
   }
   if j != Android {
      t.Fatal(j)
   }
}
