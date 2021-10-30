package binary

import (
   "encoding/base64"
   "fmt"
   "testing"
)

const androidKey = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

func TestDecode(t *testing.T) {
   data, err := base64.StdEncoding.DecodeString(androidKey)
   if err != nil {
      t.Fatal(err)
   }
   mod := NewDecoder(data).Uint32LengthPrefixed()
   fmt.Println(len(mod))
}