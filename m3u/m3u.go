// M3U parser
package m3u

import (
   "bytes"
   "io"
   "strconv"
)

func merge(forms []Format) int {
   // INSERT
   fLen := len(forms)
   if fLen == 0 {
      return -1
   }
   form := forms[fLen-1]
   if len(form) == 0 {
      return -1
   }
   // UPDATE
   return fLen-1
}

type Format map[string]string

func Decode(src io.Reader, dir string) ([]Format, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf, dir), nil
}

// #EXTINF:6.006,frame-rate=23.976
func Unmarshal(buf []byte, dir string) []Format {
   lines := bytes.FieldsFunc(buf, func(r rune) bool {
      return r == '\n'
   })
   var pass1 []Format
   for _, line := range lines {
      if line[0] == '#' {
         form := make(Format)
         pairs := reader{line}
         pairs.readBytes(':', '"')
         for {
            if pairs.buf == nil {
               break
            }
            var pair reader
            pair.buf = pairs.readBytes(',', '"')
            key := pair.readString('=', '"')
            if pair.buf != nil {
               val := string(pair.buf)
               unq, err := strconv.Unquote(val)
               if err == nil {
                  val = unq
               }
               form[key] = val
            }
         }
         pass1 = append(pass1, form)
      } else {
         text := string(line)
         ind := merge(pass1)
         if ind == -1 {
            form := Format{"URI": text}
            pass1 = append(pass1, form)
         } else {
            pass1[ind]["URI"] = text
         }
      }
   }
   var pass2 []Format
   uris := make(map[string]bool)
   for _, form := range pass1 {
      uri, ok := form["URI"]
      if ok && !uris[uri] {
         form["URI"] = dir + form["URI"]
         pass2 = append(pass2, form)
         uris[uri] = true
      }
   }
   return pass2
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

func (r *reader) readString(sep, enc byte) string {
   bytes := r.readBytes(sep, enc)
   return string(bytes)
}
