package strconv

type Ratio float64

func NewRatio[T, U Integer](value T, total U) Ratio {
   var r float64
   if total != 0 {
      r = float64(value) / float64(total)
   }
   return Ratio(r)
}

func (r Ratio) AppendPercent(b []byte) []byte {
   return label(b, r, unit{100, "%"})
}

func (r Ratio) AppendRate(b []byte) []byte {
   units := []unit{
      {1e-3, " kilobyte/s"},
      {1e-6, " megabyte/s"},
      {1e-9, " gigabyte/s"},
      {1e-12, " terabyte/s"},
   }
   return scale(b, r, units)
}
