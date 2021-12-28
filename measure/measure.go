package measure

import (
   "strconv"
)

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
)

func Percent(value, total int64) string {
   var ratio int64
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.FormatInt(ratio, 10) + " %"
}

type Symbols []string

func (s Symbols) FormatFloat(f float64) string {
   var (
      i int
      symbol string
   )
   for i, symbol = range s {
      if f < 1000 {
         break
      }
      f /= 1000
   }
   if i == 0 {
      return strconv.FormatFloat(f, 'f', 0, 64) + symbol
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
