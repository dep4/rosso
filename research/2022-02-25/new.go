package m3u

import (
   "net/url"
)

func resolve(a, b string) string {
   dir, _ := url.Parse(a)
   file, _ := url.Parse(b)
   return dir.ResolveReference(file).String()
}

func resolve2(a, b string) string {
   dir, _ := url.Parse(a)
   file, _ := dir.Parse(b)
   return file.String()
}
