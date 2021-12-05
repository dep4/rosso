package m3u

import (
   "bytes"
   "io"
   "strconv"
)

type Format struct {
   ID int
   Resolution string
   Bandwidth int
   Codecs string
   URI URI
}

type Formats []Format

func Decode(src io.Reader, dir string) (Formats, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf, dir)
}

func Unmarshal(buf []byte, dir string) (Formats, error) {
   lines := bytes.FieldsFunc(buf, func(r rune) bool {
      return r == '\n'
   })
   var pass1 Formats
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
               form.URI.File = val
            }
         }
         pass1 = append(pass1, form)
      } else {
         text := string(line)
         ind := pass1.merge()
         if ind == -1 {
            var form Format
            form.URI.File = text
            pass1 = append(pass1, form)
         } else {
            pass1[ind].URI.File = text
         }
      }
   }
   var pass2 Formats
   uris := make(map[string]bool)
   for _, form := range pass1 {
      if form.URI.File != "" && !uris[form.URI.File] {
         form.URI.Dir = dir
         form.ID = len(pass2)
         pass2 = append(pass2, form)
         uris[form.URI.File] = true
      }
   }
   return pass2, nil
}

func (f Formats) Get(n int) (Format, bool) {
   if n <= -1 {
      return Format{}, false
   }
   if n >= len(f) {
      return Format{}, false
   }
   return f[n], true
}

func (f Formats) merge() int {
   index := len(f) -1
   // INSERT
   last, ok := f.Get(index)
   if !ok || last.Resolution == "" {
      return -1
   }
   // UPDATE
   return index
}

type URI struct {
   Dir string
   File string
}

func (u URI) String() string {
   return u.Dir + u.File
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
