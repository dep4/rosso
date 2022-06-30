package hls

import (
   "strconv"
   "strings"
)

const TS = ".ts"

func (s Streams) Audio(value string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Audio, value) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) Codecs(value string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Codecs, value) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) Get_Bandwidth(value int64) *Stream {
   distance := func(s *Stream) int64 {
      if s.Bandwidth > value {
         return s.Bandwidth - value
      }
      return value - s.Bandwidth
   }
   var out *Stream
   for key, value := range s {
      if out == nil || distance(&value) < distance(out) {
         out = &s[key]
      }
   }
   return out
}

type Stream struct {
   Audio string
   Bandwidth int64
   Codecs string
   Raw_URI string
   Resolution string
}

type Streams []Stream

func (s Streams) String() string {
   var b []byte
   for i, stream := range s {
      if i >= 1 {
         b = append(b, "\n\n"...)
      }
      b = append(b, stream.String()...)
   }
   return string(b)
}

func (s Stream) String() string {
   var b []byte
   b = append(b, "Bandwidth:"...)
   b = strconv.AppendInt(b, s.Bandwidth, 10)
   if s.Codecs != "" {
      b = append(b, " Codecs:"...)
      b = append(b, s.Codecs...)
   }
   if s.Resolution != "" {
      b = append(b, " Resolution:"...)
      b = append(b, s.Resolution...)
   }
   b = append(b, "\nAudio:"...)
   b = append(b, s.Audio...)
   return string(b)
}
