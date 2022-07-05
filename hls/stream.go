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

type Stream_Func func(bandwidth int, audio, codecs string) int

func (s Streams) Streams(fn Stream_Func) Streams {
   if fn == nil {
      return s
   }
   var slice Streams
   for _, elem := range s {
      if fn(elem.Bandwidth, elem.Audio, elem.Codecs) != 0 {
         slice = append(slice, elem)
      }
   }
   return slice
}

func (s Streams) Stream(fn Stream_Func) *Stream {
   if len(s) == 0 || fn == nil {
      return nil
   }
   distance := func(i int) int {
      bandwidth, audio, codecs := s[i].Bandwidth, s[i].Audio, s[i].Codecs
      return fn(bandwidth, audio, codecs)
   }
   sort.Slice(s, func(a, b int) bool {
      return distance(a) < distance(b)
   })
   return &s[0]
}
