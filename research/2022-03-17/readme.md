I think I figured it out. You can create a wrapper type, that holds the current
object, as well as the output value. If anyone has other ideas, I am interested
in them as well:

~~~
package main

type object map[string]interface{}

type token[T any] struct {
   object
   value T
}

func newToken[T any](m object) token[T] {
   return token[T]{object: m}
}

func (t token[T]) get(key string) token[T] {
   switch val := t.object[key].(type) {
   case object:
      t.object = val
   case T:
      t.value = val
   }
   return t
}
~~~

Example:

~~~
package main

func main() {
   obj := object{
      "one": object{
         "two": object{"three": 3},
      },
   }
   three := newToken[int](obj).get("one").get("two").get("three")
   println(three.value == 3)
}
~~~
