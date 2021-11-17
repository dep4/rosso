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

type stream map[string]string

func parseStream(row string) stream {
   table := make(stream)
   cols := strings.Split(row, ",")
   for _, col := range cols {
      kv := strings.SplitN(col, "=", 2)
      if len(kv) != 2 {
         return nil
      }
      table[kv[0]] = kv[1]
   }
   return table
}
