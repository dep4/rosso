package crypto

import (
   "encoding/hex"
   "fmt"
   "os"
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

func TestFormatHandshake(t *testing.T) {
   hands := []string{androidHandshake, curlHandshake}
   for _, hand := range hands {
      data, err := hex.DecodeString(hand)
      if err != nil {
         t.Fatal(err)
      }
      hello, err := ParseHandshake(data)
      if err != nil {
         t.Fatal(err)
      }
      ja3, err := hello.FormatJA3()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(ja3)
   }
}

func TestFormatJA3(t *testing.T) {
   hello, err := ParseJA3(AndroidJA3)
   if err != nil {
      t.Fatal(err)
   }
   for _, ext := range hello.ClientHelloSpec.Extensions {
      fmt.Printf("%#v\n", ext)
   }
   ja3, err := hello.FormatJA3()
   if err != nil {
      t.Fatal(err)
   }
   if ja3 != AndroidJA3 {
      t.Fatal(ja3)
   }
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

769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0

func TestAndroid(t *testing.T) {
   hands := []byte("\x16\x03\x01\x02\x00\x01\x00\x01\xfc\x03\x03`l\xae\xe4\r\x8e\xccD\x9b\x15w{\x1cO\xc4Ř5=\xc8KM\xd5Uy\xb2\"+:g\xb4\xf5 \x9c\x19$\xdb\xd34\xe6\x11\xf0\xd7\x01\xd4\xc35w\x88\xdf-\a\n\x94\x00pN\xa2\xe1\xc3VίqK\x00\"\x13\x01\x13\x02\x13\x03\xc0+\xc0,̩\xc0/\xc00̨\xc0\t\xc0\n\xc0\x13\xc0\x14\x00\x9c\x00\x9d\x00/\x005\x01\x00\x01\x91\x00\x00\x00\x1b\x00\x19\x00\x00\x16android.googleapis.com\x00\x17\x00\x00\xff\x01\x00\x01\x00\x00\n\x00\b\x00\x06\x00\x1d\x00\x17\x00\x18\x00\v\x00\x02\x01\x00\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\r\x00\x14\x00\x12\x04\x03\b\x04\x04\x01\x05\x03\b\x05\x05\x01\b\x06\x06\x01\x02\x01uP\x00\x00\x003\x00&\x00$\x00\x1d\x00 \xa0\xe2\x99\x01\x98I<+\x83\xf6\xea<\x1d\xb4yw\xf1~\x9d\x1fmM\x0e\xe8\xef`\x90Բ%\xf8Z\x00-\x00\x02\x01\x01\x00+\x00\t\b\x03\x04\x03\x03\x03\x02\x03\x01\x00\x15\x00\xf1\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
   for _, hand := range Handshakes(hands) {
      hello, err := ParseHandshake(hand)
      if err != nil {
         t.Fatal(err)
      }
      ja3, err := hello.FormatJA3()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(ja3)
   }
}









