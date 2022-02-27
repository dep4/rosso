package m3u

import (
   "net/url"
   "path/filepath"
)

func abs(a, b string) string {
   _, err := url.Parse(b)
   if err != nil {
      if filepath.IsAbs(b) {
         return b
      }
      return filepath.Join(filepath.Dir(a), b)
   }
   return b
}
