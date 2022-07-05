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

type Streams_Func func(audio, codecs string) bool

func (s Streams) Streams(fn Streams_Func) Streams {
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

type Stream_Func func(bandwidth int) int

func (s Streams) Stream(fn Stream_Func) *Stream {
   if len(s) == 0 || fn == nil {
      return nil
   }
   distance := func(i int) int {
      bandwidth := s[i].Bandwidth
      return fn(bandwidth)
   }
   sort.Slice(s, func(a, b int) bool {
      return distance(a) < distance(b)
   })
   return &s[0]
}
