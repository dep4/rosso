package strings

import (
   "strings"
)

func CutByte(s string, sep byte) (string, string, bool) {
   i := strings.IndexByte(s, sep)
   if i == -1 {
      return s, "", false
   }
   return s[:i], s[i+1:], true
}
