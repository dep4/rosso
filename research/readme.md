type \*T is pointer to type parameter, not type parameter

I have a package like this:

~~~
package hello

type String string

func (String) name() string { return "String" }

type strings []String

func (s strings) get() *String {
   if len(s) == 0 {
      return nil
   }
   return &s[0]
}
~~~

I would like to generalize it, so I tried this:

~~~
package hello

type namer interface {
   name() string
}

type slice[T namer] []T

func (s slice[T]) get() *T {
   if len(s) == 0 {
      return nil
   }
   return &s[0]
}
~~~

but it fails:

~~~
func try_it[T namer](s slice[T]) {
   // undefined (type *T is pointer to type parameter, not type parameter)
   name := s.get().name()
   println(name)
}
~~~

how can I implement the function so that it works with different types?
