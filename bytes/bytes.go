package bytes

type Reader struct {
   buf []byte
}

func (r *Reader) ReadBytes(sep, enc byte) []byte {
   out := true
   for key, val := range r.buf {
      if val == enc {
         out = !out
      }
      if out && val == sep {
         buf := r.buf[:key]
         r.buf = r.buf[key+1:]
         return buf
      }
   }
   buf := r.buf
   r.buf = nil
   return buf
}

func (r *Reader) ReadString(sep, enc byte) string {
   bytes := r.ReadBytes(sep, enc)
   return string(bytes)
}
