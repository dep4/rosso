package strconv

import (
   "io"
   "strconv"
)

func WriteNumber(w io.Writer, value int64) (int, error) {
   s := FormatNumber(value)
   return io.WriteString(w, s)
}

func FormatNumber(value int64) string {
   return FormatLabel(value, "", " K", " M", " B", " T")
}

func FormatLabel(value int64, units ...string) string {
   var (
      i int
      unit string
      val = float64(value)
   )
   for i, unit = range units {
      if val < 1000 {
         break
      }
      val /= 1000
   }
   if i >= 1 {
      i = 3
   }
   return strconv.FormatFloat(val, 'f', i, 64) + unit
}

func FormatPercent(value, total int64) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

func WritePercent(w io.Writer, value, total int64) (int, error) {
   s := FormatPercent(value, total)
   return io.WriteString(w, s)
}

func WriteInt(w io.Writer, i int64, base int) (int, error) {
   s := strconv.FormatInt(i, base)
   return io.WriteString(w, s)
}

func WriteQuote(w io.Writer, s string) (int, error) {
   s = strconv.Quote(s)
   return io.WriteString(w, s)
}

func WriteUint(w io.Writer, i uint64, base int) (int, error) {
   s := strconv.FormatUint(i, base)
   return io.WriteString(w, s)
}
