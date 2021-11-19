package m3u

import (
   "bufio"
   "io"
   "strings"
)

func cutByte(s string, sep byte) (string, string, bool) {
   i := strings.IndexByte(s, sep)
   if i == -1 {
      return s, "", false
   }
   return s[:i], s[i+1:], true
}

// #EXT-X-BYTERANGE
type stream map[string][]string

func newStream(src io.Reader, prefix string) stream {
   str := make(stream)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         _, param, ok := cutByte(val, ':')
         if ok {
            text := prefix + buf.Text()
            params, ok := str[text]
            if ok {
               str[text] = append(params, param)
            } else {
               str[text] = []string{param}
            }
         }
      }
   }
   return str
}
