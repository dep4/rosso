package format

import (
   "strconv"
)

type Number interface {
   float64 | int | int64 | uint64
}

type Measure[T Number] struct {
   Value T
   Unit []string
}

func MeasureNumber[T Number](value T) Measure[T] {
   var m Measure[T]
   m.Value = value
   m.Unit = []string{"", " K", " M", " B", " T"}
   return m
}

func MeasureSize[T Number](value T) Measure[T] {
   var m Measure[T]
   m.Value = value
   m.Unit = []string{" B", " kB", " MB", " GB", " TB"}
   return m
}

func MeasureRate[T Number](value T) Measure[T] {
   var m Measure[T]
   m.Value = value
   m.Unit = []string{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
   return m
}

func (m Measure[T]) String() string {
   var (
      i int
      symbol string
      value = float64(m.Value)
   )
   for i, symbol = range m.Unit {
      if value < 1000 {
         break
      }
      value /= 1000
   }
   if i >= 1 {
      i = 3
   }
   return strconv.FormatFloat(value, 'f', i, 64) + symbol
}
