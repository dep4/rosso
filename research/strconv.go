package strconv

import (
   "strconv"
)

func label(value, total float64, items []scale) string {
   var ratio float64
   if total != 0 {
      ratio = value / total
   }
   var item scale
   for _, item = range items {
      if ratio < 1000 {
         break
      }
      ratio /= 1000
   }
   return strconv.FormatFloat(ratio, 'f', 3, 64) + unit
}

type scale struct {
   factor float64
   unit string
}

var number = []scale{
   {1, ""},
   {1e3, " K"},
   {1e6, " M"},
   {1e9, " B"},
   {1e12, " T"},
}

var size = []scale{
   {1, " B"},
   {1e3, " kB"},
   {1e6, " MB"},
   {1e9, " GB"},
   {1e12, " TB"},
}

var rate = []scale{
   {1, " B/s"},
   {1e3, " kB/s"},
   {1e6, " MB/s"},
   {1e9, " GB/s"},
   {1e12, " TB/s"},
}

var percent = []scale{
   {0.01, "%"},
}
