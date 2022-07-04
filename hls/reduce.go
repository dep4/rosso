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

type Reducer interface {
   Codecs() string
}

type AVC1 struct{}

func (AVC1) Codecs() string { return "avc1." }

func (s Streams) Reduce(r Reducer) Streams {
   if r == nil {
      return s
   }
   var slice Streams
   for _, elem := range s {
      if strings.Contains(elem.Codecs, r.Codecs()) {
         slice = append(slice, elem)
      }
   }
   return slice
}
