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
   URI string
}

func Formats(src io.Reader, prefix string) ([]Format, error) {
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
            dup := Format{URI: text}
            dups = append(dups, dup)
         } else {
            dups[fLen-1].URI = text
         }
      }
   }
   var forms []Format
   uris := make(map[string]bool)
   for _, dup := range dups {
      if !uris[dup.URI] {
         dup.URI = prefix + dup.URI
         forms = append(forms, dup)
         uris[dup.URI] = true
      }
   }
   return forms, nil
}

func directive(text string) (Format, error) {
   var dup Format
   str := reader{text}
   str.readString(':', '"')
   for {
      key := str.readString('=', '"')
      if key == "" {
         return dup, nil
      }
      val := str.readString(',', '"')
      switch key {
      case "BANDWIDTH":
         num, err := strconv.Atoi(val)
         if err != nil {
            return Format{}, err
         }
         dup.Bandwidth = num
      case "CODECS":
         unq, err := strconv.Unquote(val)
         if err == nil {
            val = unq
         }
         dup.Codecs = val
      case "RESOLUTION":
         dup.Resolution = val
      case "URI":
         unq, err := strconv.Unquote(val)
         if err == nil {
            val = unq
         }
         dup.URI = val
      }
   }
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
