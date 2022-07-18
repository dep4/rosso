package strconv

import (
   "strconv"
)

func Label(value float64, units []Unit) string {
   var (
      f float64
      prec int
      u Unit
   )
   for _, u = range units {
      f = value / u.Factor
      if f < 1000 {
         break
      }
   }
   if u.Factor >= 2 {
      prec = 3
   }
   return strconv.FormatFloat(f, 'f', prec, 64) + u.Name
}

func Ratio(value, total float64, units []Unit) string {
   var f float64
   if total != 0 {
      f = value / total
   }
   return Label(f, units)
}

type Unit struct {
   Factor float64
   Name string
}

var Cardinal = []Unit{
   {1, ""},
   {1e3, " thousand"},
   {1e6, " million"},
   {1e9, " billion"},
   {1e12, " trillion"},
}

var Rate = []Unit{
   {1, " byte/s"},
   {1e3, " kilobyte/s"},
   {1e6, " megabyte/s"},
   {1e9, " gigabyte/s"},
   {1e12, " terabyte/s"},
}

var Size = []Unit{
   {1, " byte"},
   {1e3, " kilobyte"},
   {1e6, " megabyte"},
   {1e9, " gigabyte"},
   {1e12, " terabyte"},
}
