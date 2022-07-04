package hls

import (
   "strings"
)

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
}

type Stream_Filter interface {
   Audio() string
   Codecs() []string
}

func (s Streams) Filter(f Stream_Filter) Streams {
   if f == nil {
      return s
   }
   pass := func(s Stream) bool {
      if !strings.Contains(s.Audio, f.Audio()) {
         return false
      }
      for _, elem := range f.Codecs() {
         if !strings.Contains(s.Codecs, elem) {
            return false
         }
      }
      return true
   }
   var slice Streams
   for _, elem := range s {
      if pass(elem) {
         slice = append(slice, elem)
      }
   }
   return slice
}
