package m3u

import (
   "bufio"
   "io"
   "strings"
)

const stream = "#EXT-X-STREAM-INF:"

type playlist map[string]string

func newPlaylist(r io.Reader) playlist {
   list := make(playlist)
   buf := bufio.NewScanner(r)
   for buf.Scan() {
      key := buf.Text()
      if strings.HasPrefix(key, stream) {
         buf.Scan()
         list[strings.TrimPrefix(key, stream)] = buf.Text()
      }
   }
   return list
}
