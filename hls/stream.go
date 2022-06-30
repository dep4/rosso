package hls

import (
   "bytes"
   "strconv"
   "strings"
)

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
   Resolution string
   URI string
}

type Streams []Stream

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
   b = append(b, "\n\tAudio:"...)
   b = append(b, s.Audio...)
   return string(b)
}

func (s Stream) Ext(b []byte) string {
   if bytes.Contains(b, []byte("ftypiso5")) {
      return ".m4v"
   }
   if bytes.HasPrefix(b, []byte{'G'}) {
      return ".ts"
   }
   return ""
}
