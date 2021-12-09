package crypto

import (
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

const androidHandshake =
   "16030100bb010000b703034420d198e7852decbc117dc7f90550b98f2d643c954bf3361d" +
   "daf127ff921b04000024c02bc02ccca9c02fc030cca8009e009fc009c00ac013c0140033" +
   "0039009c009d002f00350100006aff0100010000000022002000001d636c69656e747365" +
   "7276696365732e676f6f676c65617069732e636f6d0017000000230000000d0016001406" +
   "010603050105030401040303010303020102030010000b000908687474702f312e31000b" +
   "00020100000a000400020017"

const androidJA3 =
   "769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-" +
   "51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

func TestHandshake(t *testing.T) {
   data, err := hex.DecodeString(androidHandshake)
   if err != nil {
      t.Fatal(err)
   }
   hello, err := ParseHandshake(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", hello)
}

func TestHandshakes(t *testing.T) {
   pcap, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      t.Fatal(err)
   }
   for _, hand := range Handshakes(pcap) {
      hello, err := ParseHandshake(hand)
      if err == nil {
         fmt.Printf("%+v\n", hello)
      }
   }
}

func TestJA3(t *testing.T) {
   h, err := ParseJA3(androidJA3)
   if err != nil {
      t.Fatal(err)
   }
   for _, ext := range h.ClientHelloSpec.Extensions {
      fmt.Printf("%#v\n", ext)
   }
   j, err := h.FormatJA3()
   if err != nil {
      t.Fatal(err)
   }
   if j != androidJA3 {
      t.Fatal(j)
   }
}
