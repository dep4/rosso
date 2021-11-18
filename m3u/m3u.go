package m3u

import (
   "bufio"
   "io"
   "strings"
)

func Parse(src io.Reader, prefix string) map[string][]string {
   list := make(map[string][]string)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         _, param, ok := cutByte(val, ':')
         if ok {
            text := prefix + buf.Text()
            params, ok := list[text]
            if ok {
               list[text] = append(params, param)
            } else {
               list[text] = []string{param}
            }
         }
      }
   }
   return list
}

func cutByte(s string, sep byte) (string, string, bool) {
   i := strings.IndexByte(s, sep)
   if i == -1 {
      return s, "", false
   }
   return s[:i], s[i+1:], true
}
