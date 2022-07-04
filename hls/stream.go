package hls

import (
   "strconv"
   "strings"
)

func (s Streams) Video() Streams {
   var slice Streams
   for _, elem := range s {
      if elem.Codecs != "" && elem.Resolution == "" {
         continue
      }
      slice = append(slice, elem)
   }
   return slice
}

func (s Streams) Bandwidth(value int) *Stream {
   distance := func(s *Stream) int {
      if s.Bandwidth > value {
         return s.Bandwidth - value
      }
      return value - s.Bandwidth
   }
   var elem *Stream
   for key, value := range s {
      if elem == nil || distance(&value) < distance(elem) {
         elem = &s[key]
      }
   }
   return elem
}

func (s Stream) String() string {
   var (
      a []string
      b string
   )
   if s.Resolution != "" {
      a = append(a, "Resolution:" + s.Resolution)
   }
   a = append(a, "Bandwidth:" + strconv.Itoa(s.Bandwidth))
   if s.Codecs != "" {
      a = append(a, "Codecs:" + s.Codecs)
   }
   if s.Audio != "" {
      b = "Audio:" + s.Audio
   }
   ja := strings.Join(a, " ")
   if b != "" {
      return ja + "\n\t" + b
   }
   return ja
}

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
}

type Streams []Stream

func (Stream) Ext() string {
   return ".m4v"
}
