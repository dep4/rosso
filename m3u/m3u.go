package m3u

import (
   "bufio"
   "io"
   "strconv"
   "strings"
)

type Playlist map[string]map[string]string

func NewPlaylist(src io.Reader, prefix string) Playlist {
   list := make(Playlist)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         str := reader{val}
         str.readString(':', '"')
         param := make(map[string]string)
         for {
            key := str.readString('=', '"')
            if key == "" {
               break
            }
            val := str.readString(',', '"')
            unq, err := strconv.Unquote(val)
            if err != nil {
               param[key] = val
            } else {
               param[key] = unq
            }
         }
         list[prefix + buf.Text()] = param
      }
   }
   return list
}

type Stream map[string][]string

func NewStream(src io.Reader, prefix string) Stream {
   str := make(Stream)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         param := reader{val}
         param.readString(':', '"')
         text := prefix + buf.Text()
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
