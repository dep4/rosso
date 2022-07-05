package hls

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
}

type Stream_Filter func(Stream) bool

func (s Streams) Filter(callback Stream_Filter) Streams {
   var carry Streams
   for _, item := range s {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

type Stream_Reduce func(*Stream, Stream) *Stream

func (s Streams) Reduce(callback Stream_Reduce) *Stream {
   var carry *Stream
   for _, item := range s {
      carry = callback(carry, item)
   }
   return carry
}
