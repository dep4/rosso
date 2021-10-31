package binary

import (
   "bytes"
   "encoding/binary"
   "github.com/89z/parse/tls"
   "io"
)

func handshake(data []byte) *tls.ClientHello {
   r := bytes.NewReader(data)
   for {
      for {
         typ, err := r.ReadByte()
         if err != nil {
            return nil
         }
         if typ == 0x16 {
            break
         }
      }
      w := new(bytes.Buffer)
      w.WriteByte(0x16)
      io.CopyN(w, r, 2)
      buf := make([]byte, 2)
      r.Read(buf)
      w.Write(buf)
      off := int64(binary.BigEndian.Uint16(buf))
      _, err := io.CopyN(w, r, off)
      if err == nil {
         hello, err := tls.ParseHandshake(w.Bytes())
         if err == nil {
            return hello
         }
      }
      r.Seek(-off, io.SeekCurrent)
   }
}
