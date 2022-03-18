package world
import "hello"

func IsZero[T int64|float64](value T) bool {
   return hello.IsZero(value)
}
