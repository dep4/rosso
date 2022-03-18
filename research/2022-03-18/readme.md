If the function is small, like in the question, its probably easier to just
vendor it:

~~~go
package vendor

func thisIsJustCopy[T int64|float64](value T) bool {
   return value == 0
}
~~~

but if the function is big, you can do it like this:

~~~go
package world
import "hello"

func IsZero[T int64|float64](value T) bool {
   return hello.IsZero(value)
}
~~~
