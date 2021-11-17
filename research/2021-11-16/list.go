package m3u

import (
   "bufio"
   "io"
   "strings"
)

type playlist map[string][]string

func newPlaylist(r io.Reader) playlist {
   list := make(playlist)
   buf := bufio.NewScanner(r)
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), ":", 2)
      if strings.HasPrefix(key, "stream") {
         buf.Scan()
         list[strings.TrimPrefix(key, "stream")] = buf.Text()
      }
   }
   return list
}
