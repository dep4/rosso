package crypto

import (
   "fmt"
   "net/http"
   "testing"
)

func TestFormatJA3(t *testing.T) {
   hello, err := ParseJA3(AndroidAPI26)
   if err != nil {
      t.Fatal(err)
   }
   ja3, err := FormatJA3(hello)
   if err != nil {
      t.Fatal(err)
   }
   if ja3 != AndroidAPI26 {
      t.Fatal(ja3)
   }
}

func TestTransport(t *testing.T) {
   req, err := http.NewRequest("HEAD", "https://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   hello, err := ParseJA3(AndroidAPI26)
   if err != nil {
      t.Fatal(err)
   }
   res, err := Transport(hello).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", res)
}

func TestReader(t *testing.T) {
   var buf []byte
   buf = append(buf, 0,0,0,5, 'h', 'e', 'l', 'l', 'o')
   buf = append(buf, 0,0,0,5, 'w', 'o', 'r', 'l', 'd')
   read := NewReader(buf)
   one, two, ok := read.ReadUint32LengthPrefixed()
   fmt.Printf("%v %s %v\n", one, two, ok)
   one, two, ok = read.ReadUint32LengthPrefixed()
   fmt.Printf("%v %s %v\n", one, two, ok)
}
