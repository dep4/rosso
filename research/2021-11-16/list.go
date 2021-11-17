package m3u

import (
   "bufio"
   "io"
   "strings"
)

type playlist map[string]string

func newPlaylist(r io.Reader) playlist {
   list := make(playlist)
   buf := bufio.NewScanner(r)
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), ":", 2)
      if strings.HasPrefix(kv[0], "stream") {
         buf.Scan()
         list[strings.TrimPrefix(buf.Text(), "stream")] = buf.Text()
      }
   }
   return list
}
