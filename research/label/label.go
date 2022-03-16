package format

import (
   "strconv"
)

type number interface {
   float64 | int | int64 | uint64
}

type measure[T number] struct {
   unit []string
}

func (m measure[T]) label(value T) string {
   var (
      i int
      symbol string
      val = float64(value)
   )
   for i, symbol = range m.unit {
      if val < 1000 {
         break
      }
      val /= 1000
   }
   if i >= 1 {
      i = 3
   }
   return strconv.FormatFloat(val, 'f', i, 64) + symbol
}

func (m *measure[T]) labelNumber(value T) string {
   m.unit = []string{"", " K", " M", " B", " T"}
   return m.label(value)
}

func (m *measure[T]) labelRate(value T) string {
   m.unit = []string{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
   return m.label(value)
}

func (m *measure[T]) labelSize(value T) string {
   m.unit = []string{" B", " kB", " MB", " GB", " TB"}
   return m.label(value)
}
