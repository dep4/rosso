package binary

import (
   "encoding/base64"
   "fmt"
   "os"
   "testing"
)

const androidKey = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

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
   hand := Handshake(data)
   fmt.Printf("Handshake %+v\n", hand)
   hand = handshake(data)
   fmt.Printf("handshake %+v\n", hand)
}
