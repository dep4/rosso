package binary

import (
   "bytes"
   "github.com/89z/parse/binary"
   "github.com/89z/parse/tls"
)

func Handshake(b []byte) *tls.ClientHello {
   for {
      rec1 := bytes.IndexByte(b, 0x16)
      if rec1 == -1 {
         return nil
      }
      ver1 := rec1 + 1
      len1 := ver1 + 2
      recLen, ok := binary.Uint16(b[len1:])
      if ok {
         len2 := len1 + 2
         rec2 := len2 + int(recLen)
         if rec2 < len(b) {
            hand, err := tls.ParseHandshake(b[rec1:rec2])
            if err == nil {
               return hand
            }
         }
      }
      b = b[rec1+1:]
   }
}
