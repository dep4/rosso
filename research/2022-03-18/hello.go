package hello

func IsZero[T int64|float64](value T) bool {
   return value == 0
}
