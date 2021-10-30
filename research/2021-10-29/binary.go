package binary

import (
   "bytes"
   "github.com/89z/parse/binary"
   "io"
)

func handshakes(data []byte) [][]byte {
   r := bytes.NewReader(data)
   var hands [][]byte
   for {
      for {
         typ, err := r.ReadByte()
         if err != nil {
            return hands
         }
         if typ == 0x16 {
            break
         }
      }
      w := new(bytes.Buffer)
      // Content Type
      io.CopyN(w, r, 1)
      // Version
      io.CopyN(w, r, 2)
      // Length, Handshake Protocol
      var buf [2]byte
      r.Read(buf[:])
      w.Write(buf[:])
      off, err := io.CopyN(w, r, binary.Varint(buf[:]))
      if err != nil {
         r.Seek(-off, io.SeekCurrent)
      } else {
         hands = append(hands, w.Bytes())
      }
   }
}
