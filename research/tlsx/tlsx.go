package main

import (
   "fmt"
   "github.com/89z/format/crypto"
   "github.com/dreadl0ck/ja3"
   "github.com/dreadl0ck/tlsx"
)

// this is the capture from Android API 29. For some reason, the JA3 below does
// not match mine.
var payload = []byte("\x16\x03\x01\x02\x00\x01\x00\x01\xfc\x03\x03r-\x02>\xc9\x1a\xa3\x93\xc3\xeb9\x10\xed\xc4U\xb3\xd1\\/\x040F\x14\xd7Z\x1f5&\x9cXX\x81 \xecw\x9e\x8cl\r\x9f\xacg[\x81\xb6Y\xed\xb1\xe2l\xe6 G\x93\n$~w\x95z\xb5<\a\x10\x8a\x00\"\x13\x01\x13\x02\x13\x03\xc0+\xc0,̩\xc0/\xc00̨\xc0\t\xc0\n\xc0\x13\xc0\x14\x00\x9c\x00\x9d\x00/\x005\x01\x00\x01\x91\x00\x00\x00\"\x00 \x00\x00\x1dclientservices.googleapis.com\x00\x17\x00\x00\xff\x01\x00\x01\x00\x00\n\x00\b\x00\x06\x00\x1d\x00\x17\x00\x18\x00\v\x00\x02\x01\x00\x00#\x00\x00\x00\x10\x00\v\x00\t\bhttp/1.1\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\r\x00\x14\x00\x12\x04\x03\b\x04\x04\x01\x05\x03\b\x05\x05\x01\b\x06\x06\x01\x02\x01\x003\x00&\x00$\x00\x1d\x00 \xc1\x1c9\x1dӜX\xe1\u007f\xc0o+l\xa4|\xbed\xee\x96\x19/\xb3\xa9\x1e8\xda\xe7O~a\xc1\x1c\x00-\x00\x02\x01\x01\x00+\x00\t\b\x03\x04\x03\x03\x03\x02\x03\x01\x00\x15\x00\xdb\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")

func main() {
   {
      var hello tlsx.ClientHelloBasic
      err := hello.Unmarshal(payload)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%s\n", ja3.Bare(&hello))
   }
   {
      hello, err := crypto.ParseTLS(payload)
      if err != nil {
         panic(err)
      }
      ja3, err := crypto.FormatJA3(hello)
      if err != nil {
         panic(err)
      }
      fmt.Println(ja3)
   }
}
