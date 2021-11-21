package m3u

import (
   "bufio"
   "io"
   "strconv"
   "strings"
)

type ByteRange map[string][]string

func NewByteRange(src io.Reader) ByteRange {
   str := make(ByteRange)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if strings.HasPrefix(text, "#") {
         val = text
      } else {
         param := reader{val}
         param.readString(':', '"')
         params, ok := str[text]
         if ok {
            str[text] = append(params, param.str)
         } else {
            str[text] = []string{param.str}
         }
      }
   }
   return str
}

type Directive map[string]string

func newDirective(src string) Directive {
   str := reader{src}
   str.readString(':', '"')
   dir := make(Directive)
   for {
      key := str.readString('=', '"')
      if key == "" {
         return dir
      }
      val := str.readString(',', '"')
      unq, err := strconv.Unquote(val)
      if err == nil {
         val = unq
      }
      dir[key] = val
   }
}

type Playlist map[string]Directive

func NewPlaylist(src io.Reader) Playlist {
   list := make(Playlist)
   var val Directive
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if strings.HasPrefix(text, "#") {
         dir := newDirective(text)
         uri, ok := dir["URI"]
         if ok {
            delete(dir, "URI")
            list[uri] = dir
         } else {
            val = dir
         }
      } else {
         list[text] = val
      }
   }
   return list
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
