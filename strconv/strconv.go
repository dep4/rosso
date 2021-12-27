package strconv

import (
   "strconv"
)

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
)

type Symbols []string

func (s Symbols) FormatFloat(f float64) string {
   var symbol string
   for _, symbol = range s {
      if f < 1000 {
         break
      }
      f /= 1000
   }
   return strconv.FormatFloat(f, 'f', 3, 64) + symbol
}

func (s Symbols) FormatInt(i int64) string {
   f := float64(i)
   return s.FormatFloat(f)
}

func (s Symbols) FormatUint(i uint64) string {
   f := float64(i)
   return s.FormatFloat(f)
}
