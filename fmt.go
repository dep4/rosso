package format

import (
   "io"
   "fmt"
)

func percent(w io.Writer, value, total float64) (int, error) {
   var ratio float64
   if total != 0 {
      ratio = 100 * value / total
   }
   return fmt.Fprintf(w, "%.1v%%", ratio)
}
