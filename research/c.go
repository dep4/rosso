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

func try_it[T namer](s slice[T]) {
   name := s.get().name()
   println(name)
}
