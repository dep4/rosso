package m3u

import (
   "bytes"
   "io"
   "strconv"
)

func merge(forms []Format) int {
   fLen := len(forms)
   if fLen == 0 {
      return -1
   }
   form := forms[fLen-1]
   if form.Resolution == "" {
      // INSERT
      return -1
   }
   // UPDATE
   return fLen-1
}

type Format struct {
   ID int
   Resolution string
   Bandwidth int
   Codecs string
   URI string
}

func Decode(src io.Reader) ([]Format, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

func Unmarshal(buf []byte) ([]Format, error) {
   lines := bytes.FieldsFunc(buf, func(r rune) bool {
      return r == '\n'
   })
   var pass1 []Format
   for _, line := range lines {
      if line[0] == '#' {
         var form Format
         com := reader{line}
         com.readBytes(':', '"')
         for {
            key := com.readString('=', '"')
            val := com.readString(',', '"')
            if val == "" {
               break
            }
            switch key {
            case "BANDWIDTH":
               num, err := strconv.Atoi(val)
               if err != nil {
                  return nil, err
               }
               form.Bandwidth = num
            case "CODECS":
               unq, err := strconv.Unquote(val)
               if err == nil {
                  val = unq
               }
               form.Codecs = val
            case "RESOLUTION":
               form.Resolution = val
            case "URI":
               unq, err := strconv.Unquote(val)
               if err == nil {
                  val = unq
               }
               form.URI = val
            }
         }
         pass1 = append(pass1, form)
      } else {
         text := string(line)
         ind := merge(pass1)
         if ind == -1 {
            form := Format{URI: text}
            pass1 = append(pass1, form)
         } else {
            pass1[ind].URI = text
         }
      }
   }
   var pass2 []Format
   uris := make(map[string]bool)
   for _, form := range pass1 {
      if form.URI != "" && !uris[form.URI] {
         form.ID = len(pass2)
         pass2 = append(pass2, form)
         uris[form.URI] = true
      }
   }
   return pass2, nil
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
