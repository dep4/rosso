package format

import (
   "fmt"
)

func AlfaLabel[T int|int64](value T, unit string) string {
   return fmt.Sprint(value, " ", unit)
}

type Bravo[T int|int64] struct {
   Value T
}

func (b Bravo[T]) Label(unit string) string {
   return fmt.Sprint(b.Value, " ", unit)
}

type Charlie[T int|int64] struct {
   Unit string
}

func (c Charlie[T]) Label(value T) string {
   return fmt.Sprint(value, " ", c.Unit)
}
