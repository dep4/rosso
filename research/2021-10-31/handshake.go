package bytes

import (
   "github.com/89z/parse/tls"
)

func handshake(data []byte) *tls.ClientHello {
   for {
      var w []byte
      r := buffer{data}
      low, ok := r.readBytes(0x16)
      if ! ok {
         return nil
      }
      w = append(w, 0x16)
      buf, ok := r.next(2)
      if ok {
         w = append(w, buf...)
      }
      pre, buf, ok := r.readUint16LengthPrefixed()
      if ok {
         w = append(w, pre...)
         w = append(w, buf...)
      }
      hand, err := tls.ParseHandshake(w)
      if err == nil {
         return hand
      }
      data = data[len(low):]
   }
}
