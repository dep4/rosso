package m3u

import (
   "net/url"
   "path/filepath"
)

func undecorate(v *url.URL) string {
   return filepath.FromSlash(v.Path)
}

func decorate(v string) *url.URL {
   var a url.URL
   a.Scheme = "file"
   a.Path = filepath.ToSlash(v)
   return &a
}
