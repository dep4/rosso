package hls

import (
   "strconv"
   "strings"
)

func (s Stream) String() string {
   var b []byte
   b = append(b, "Bandwidth:"...)
   b = strconv.AppendInt(b, s.Bandwidth, 10)
   if s.Raw_Codecs != "" {
      b = append(b, " Codecs:"...)
      b = append(b, s.Codecs()...)
   }
   if s.Resolution != "" {
      b = append(b, " Resolution:"...)
      b = append(b, s.Resolution...)
   }
   b = append(b, " Range:"...)
   b = append(b, s.Video_Range...)
   return string(b)
}

func (s Stream) Codecs() string {
   codecs := strings.Split(s.Raw_Codecs, ",")
   for i, codec := range codecs {
      before, _, found := strings.Cut(codec, ".")
      if found {
         codecs[i] = before
      }
   }
   return strings.Join(codecs, ",")
}

type Streams []Stream

type Stream struct {
   Bandwidth int64
   Raw_Codecs string
   Raw_URI string
   Resolution string
   Video_Range string
}

// use AUDIO instead
func (s Streams) URI(value string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Raw_URI, value) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) Codecs(value string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Raw_Codecs, value) {
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

func (s Streams) Video_Range(value string) Streams {
   var out Streams
   for _, stream := range s {
      if stream.Video_Range == value {
         out = append(out, stream)
      }
   }
   return out
}
