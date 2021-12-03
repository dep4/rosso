package m3u

import (
   "bufio"
   "io"
   "strconv"
   "strings"
)

func isDirective(text string) bool {
   return strings.HasPrefix(text, "#EXT-X-STREAM-INF:")
}

func isURI(text string) bool {
   return text != "" && text[0] != '#'
}

type Format struct {
   ID int
   Resolution string
   Bandwidth int
   Codecs string
   URI URI
}

func Formats(src io.Reader, dir string) ([]Format, error) {
   var dups []Format
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if isDirective(text) {
         dup, err := directive(text)
         if err != nil {
            return nil, err
         }
         dup.ID = len(dups)
         dups = append(dups, dup)
      }
      if isURI(text) {
         fLen := len(dups)
         if fLen == 0 {
            dups = append(dups, Format{
               URI: URI{File: text},
            })
         } else {
            dups[fLen-1].URI.File = text
         }
      }
   }
   var forms []Format
   uris := make(map[string]bool)
   for _, dup := range dups {
      if !uris[dup.URI.File] {
         dup.URI.Dir = dir
         forms = append(forms, dup)
         uris[dup.URI.File] = true
      }
   }
   return forms, nil
}

func directive(text string) (Format, error) {
   var form Format
   str := reader{text}
   str.readString(':', '"')
   for {
      key := str.readString('=', '"')
      if key == "" {
         return form, nil
      }
      val := str.readString(',', '"')
      switch key {
      case "BANDWIDTH":
         num, err := strconv.Atoi(val)
         if err != nil {
            return Format{}, err
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
}

type URI struct {
   Dir string
   File string
}

func (u URI) String() string {
   return u.Dir + u.File
}

type reader struct {
   str string
}

func (r *reader) readString(sep, enc rune) string {
   out := true
   for k, v := range r.str {
      if v == enc {
         out = !out
      }
      if out && v == sep {
         str := r.str[:k]
         r.str = r.str[k+1:]
         return str
      }
   }
   str := r.str
   r.str = ""
   return str
}
