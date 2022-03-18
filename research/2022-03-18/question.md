Alias generic function

I can define a generic function:

~~~go
package hello

func IsZero[T int64|float64](value T) bool {
   return value == 0
}
~~~

but if I try to alias that function in another package, it fails:

~~~go
package world
import "hello"

// cannot use generic function hello.IsZero without instantiation
var IsZero = hello.IsZero

// this works
// var IsZero = hello.IsZero[int64]
~~~

Is it possible to do this, using some other syntax?
