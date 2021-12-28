package format

import (
   "strconv"
)

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
)

// godocs.io/github.com/google/pprof/internal/measurement#Percentage
func Percent(value, total int) string {
   var ratio int
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.Itoa(ratio) + "%"
}

// godocs.io/github.com/google/pprof/internal/measurement#Percentage
func PercentInt64(value, total int64) string {
   var ratio int64
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.FormatInt(ratio, 10) + "%"
}

// github.com/golang/text/blob/18b340fc/encoding/internal/enctest/enctest.go#L175-L180
func Trim(s string) string {
   if len(s) <= 99 {
      return s
   }
   return s[:48] + "..." + s[len(s)-48:]
}

type Symbols []string

// godocs.io/github.com/google/pprof/internal/measurement#Label
func (s Symbols) Label(f float64) string {
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

// godocs.io/github.com/google/pprof/internal/measurement#Label
func (s Symbols) LabelInt(i int64) string {
   f := float64(i)
   return s.Label(f)
}

// godocs.io/github.com/google/pprof/internal/measurement#Label
func (s Symbols) LabelUint(i uint64) string {
   f := float64(i)
   return s.Label(f)
}
