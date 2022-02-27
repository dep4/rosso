package m3u

import (
   "net/url"
   "path"
   "path/filepath"
)

func abs(dir, file string) string {
   _, err := url.ParseRequestURI(file)
   if err != nil {
      addr, err := url.ParseRequestURI(dir)
      if err != nil || addr.Host == "" {
         return filepath.Join(filepath.Dir(dir), file)
      }
      addr.Path = path.Join(path.Dir(addr.Path), file)
      return addr.String()
   }
   return file
}
