package m3u

import (
   "bufio"
   "io"
   "strings"
)

type format struct {
   resolution, uri string
}

func decode(src io.Reader) []format {
   var forms []format
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if strings.HasPrefix(text, "#EXT-X-STREAM-INF:") {
      }
   }
   return forms
}

type reader struct {
   buf []byte
}

func (r *reader) readBytes(sep, enc byte) []byte {
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
