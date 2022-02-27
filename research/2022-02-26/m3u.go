package m3u

import (
   "net/url"
   "path/filepath"
)

const (
   specAbsolute = iota
   specRelative
   specURL
)

func specification(entry string) int {
   if filepath.IsAbs(entry) {
      return specAbsolute
   }
   _, err := url.ParseRequestURI(entry)
   if err != nil {
      return specRelative
   }
   return specURL
}
