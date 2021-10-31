package binary

import (
   "bytes"
   "github.com/89z/parse/binary"
   "github.com/89z/parse/tls"
)

func Handshake(data []byte) *tls.ClientHello {
   for {
      rec1 := bytes.IndexByte(data, 0x16)
      if rec1 == -1 {
         return nil
      }
      ver1 := rec1 + 1
      len1 := ver1 + 2
      recLen, ok := binary.Uint16(data[len1:])
      if ok {
         len2 := len1 + 2
         rec2 := len2 + int(recLen)
         if rec2 < len(data) {
            hello, err := tls.ParseHandshake(data[rec1:rec2])
            if err == nil {
               return hello
            }
         }
      }
      data = data[rec1+1:]
   }
}
