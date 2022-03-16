package format

import (
   "strconv"
)

func PercentInt(value, total int) string {
   val, tot := float64(value), float64(total)
   return percent(val, tot)
}

func PercentInt64(value, total int64) string {
   val, tot := float64(value), float64(total)
   return percent(val, tot)
}

func percent(value, total float64) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}
