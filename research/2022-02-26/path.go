package m3u

import (
   "net/url"
   "path/filepath"
)

func isAbsPath(path string) bool {
   addr, err := url.Parse(path)
   if err != nil {
      return filepath.IsAbs(path)
   }
   return addr.IsAbs()
}
