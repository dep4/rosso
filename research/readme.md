Restrict types with generic function

Currently I have a type like this:

~~~go
package hello

type object map[string]any

func add[T any](obj object, key string, val T) {
   switch value := obj[key].(type) {
   case nil:
      obj[key] = val
   case T:
      obj[key] = []T{value, val}
   case []T:
      obj[key] = append(value, val)
   }
}
~~~

which I use to store different types of numbers. I was thinking about
restricting the allowed types, like this:

~~~go
package hello

type number interface {
   isNumber()
}

type Int32 int32
func (Int32) isNumber(){}

type Float32 float32
func (Float32) isNumber(){}

type object map[string]number
~~~

but I am not sure how to implement my add function as before. I tried the same
function, but I get this:

~~~
cannot use val (variable of type T constrained by any) as type number in assignment:
   T does not implement number (missing isNumber method)
~~~

so then I changed the signature to:

~~~go
func add[T number](obj object, key string, val T)
~~~

but I get another error:

~~~
cannot use []T{…} (value of type []T) as type number in assignment:
   []T does not implement number (missing isNumber method)
~~~

is it possible to do something like what I am trying to do?
