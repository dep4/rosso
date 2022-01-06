package crypto

import (
   "encoding/hex"
   "fmt"
   "net/http"
   "testing"
)

func TestBytes(t *testing.T) {
   var b []byte
   b = append(b, 0,0,0,5, 'h', 'e', 'l', 'l', 'o')
   b = append(b, 0,0,0,5, 'w', 'o', 'r', 'l', 'd')
   buf := NewBuffer(b)
   one, two, ok := buf.ReadUint32LengthPrefixed()
   fmt.Printf("%v %s %v\n", one, two, ok)
   one, two, ok = buf.ReadUint32LengthPrefixed()
   fmt.Printf("%v %s %v\n", one, two, ok)
}

const androidHandshake =
   "16030100bb010000b703034420d198e7852decbc117dc7f90550b98f2d643c954bf3361d" +
   "daf127ff921b04000024c02bc02ccca9c02fc030cca8009e009fc009c00ac013c0140033" +
   "0039009c009d002f00350100006aff0100010000000022002000001d636c69656e747365" +
   "7276696365732e676f6f676c65617069732e636f6d0017000000230000000d0016001406" +
   "010603050105030401040303010303020102030010000b000908687474702f312e31000b" +
   "00020100000a000400020017"

const curlHandshake =
   "1603010200010001fc03033356ee099c006213ecb9f7493ef981dd513761eae27eff36a1" +
   "77ebd353fc207520fa9ef53871b81af022e38d46ca9268be95889d6e964db818768ec86a" +
   "68c7216f003e130213031301c02cc030009fcca9cca8ccaac02bc02f009ec024c028006b" +
   "c023c0270067c00ac0140039c009c0130033009d009c003d003c0035002f00ff01000175" +
   "00000010000e00000b6578616d706c652e636f6d000b000403000102000a000c000a001d" +
   "0017001e00190018337400000010000e000c02683208687474702f312e31001600000017" +
   "000000310000000d0030002e040305030603080708080809080a080b0804080508060401" +
   "05010601030302030301020103020202040205020602002b000908030403030302030100" +
   "2d00020101003300260024001d002034107e2fb61cbfc3c827b3d574b57d9d5f5294bedb" +
   "7ee350407c05d1a9396b5b001500b2000000000000000000000000000000000000000000" +
   "000000000000000000000000000000000000000000000000000000000000000000000000" +
   "000000000000000000000000000000000000000000000000000000000000000000000000" +
   "000000000000000000000000000000000000000000000000000000000000000000000000" +
   "000000000000000000000000000000000000000000000000000000000000000000000000" +
   "00000000000000000000000000"

func TestHandshake(t *testing.T) {
   hands := []string{androidHandshake, curlHandshake}
   for _, hand := range hands {
      data, err := hex.DecodeString(hand)
      if err != nil {
         t.Fatal(err)
      }
      hello, err := ParseTLS(data)
      if err != nil {
         t.Fatal(err)
      }
      ja3, err := FormatJA3(hello)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(ja3)
   }
}

func TestJA3(t *testing.T) {
   hello, err := ParseJA3(AndroidJA3)
   if err != nil {
      t.Fatal(err)
   }
   for _, ext := range hello.Extensions {
      fmt.Printf("%#v\n", ext)
   }
   ja3, err := FormatJA3(hello)
   if err != nil {
      t.Fatal(err)
   }
   if ja3 != AndroidJA3 {
      t.Fatal(ja3)
   }
}

func TestTransport(t *testing.T) {
   hello, err := ParseJA3(AndroidJA3)
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest("HEAD", "https://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   res, err := Transport(hello).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", res)
}
