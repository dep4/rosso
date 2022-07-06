package hls

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

type Element interface {
   Medium | Stream
}

type Filter[T Element] func(T) bool

type Reduce[T Element] func(*T, T) *T

func filter[T Element](array []T, callback Filter[T]) []T {
   if callback == nil {
      return array
   }
   var carry []T
   for _, item := range array {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (m Media) Filter(callback Filter[Medium]) Media {
   return filter(m, callback)
}

func reduce[T Element](array []T, callback Reduce[T]) *T {
   if callback == nil {
      return nil
   }
   var carry *T
   for _, item := range array {
      carry = callback(carry, item)
   }
   return carry
}

func (m Media) Reduce(callback Reduce[Medium]) *Medium {
   return reduce(m, callback)
}

