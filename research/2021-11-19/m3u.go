package m3u

type buffer struct {
   buf []byte
}

func newBuffer(buf []byte) *buffer {
   return &buffer{buf}
}

func (b *buffer) readBytes(sep, enc byte) []byte {
   // remove some buf
   out := true
   for k, v := range b.buf {
      if v == enc {
         out = !out
      }
      if out && v == sep {
         buf := b.buf[:k]
         b.buf = b.buf[k+1:]
         return buf
      }
   }
   // remove all buf
   buf := b.buf
   b.buf = nil
   return buf
}
