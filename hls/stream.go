package hls

import (
   "sort"
)

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
}

type Stream_Filter func(audio, codecs string) bool

func (s Streams) Filter(fn Stream_Filter) Streams {
   if fn == nil {
      return s
   }
   var slice Streams
   for _, elem := range s {
      if fn(elem.Audio, elem.Codecs) {
         slice = append(slice, elem)
      }
   }
   return slice
}

type Bandwidth func(int) int

func (s Streams) Reduce(fn Bandwidth) *Stream {
   if len(s) == 0 {
      return nil
   }
   sort.Slice(s, func(a, b int) bool {
      sa, sb := s[a].Bandwidth, s[b].Bandwidth
      return fn(sa) < fn(sb)
   })
   return &s[0]
}
