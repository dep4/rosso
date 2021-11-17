package m3u

import (
   "strings"
)

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
