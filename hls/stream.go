package hls

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
}

func (s Streams) Filter(callback Filter[Stream]) Streams {
   if callback == nil {
      return s
   }
   var carry Streams
   for _, item := range s {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (s Streams) Reduce(callback Reduce[Stream]) *Stream {
   if callback == nil {
      return nil
   }
   var carry *Stream
   for _, item := range s {
      carry = callback(carry, item)
   }
   return carry
}

func Bandwidth(v int) Reduce[Stream] {
   distance := func(s *Stream) int {
      if s.Bandwidth > v {
         return s.Bandwidth - v
      }
      return v - s.Bandwidth
   }
   return func(carry *Stream, item Stream) *Stream {
      if carry == nil || distance(&item) < distance(carry) {
         return &item
      }
      return carry
   }
}
