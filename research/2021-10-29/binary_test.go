package binary

import (
   "encoding/base64"
   "fmt"
   "github.com/89z/parse/tls"
   "os"
   "testing"
)

const androidKey = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

func TestVarint(t *testing.T) {
   n := varint([]byte{1})
   fmt.Println(n)
   n = varint([]byte{1, 0})
   fmt.Println(n)
   n = varint([]byte{1, 0, 0})
   fmt.Println(n)
   n = varint([]byte{1, 0, 0, 0})
   fmt.Println(n)
}

func TestDecode(t *testing.T) {
   _, err := base64.StdEncoding.DecodeString(androidKey)
   if err != nil {
      t.Fatal(err)
   }
}

func TestHandshakes(t *testing.T) {
   data, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      t.Fatal(err)
   }
   hands, err := handshakes(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, hand := range hands {
      hello, err := tls.ParseHandshake(hand)
      if err == nil {
         fmt.Printf("%+v\n", hello)
      }
   }
}
