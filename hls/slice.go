package hls

import (
   "strconv"
   "strings"
)

func (Medium) Ext() string {
   return ".m4a"
}

func (m Medium) String() string {
   var b strings.Builder
   b.WriteString("Type:")
   b.WriteString(m.Type)
   b.WriteString(" Name:")
   b.WriteString(m.Name)
   b.WriteString("\n  Group ID:")
   b.WriteString(m.Group_ID)
   if m.Characteristics != "" {
      b.WriteString("\n  Characteristics:")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

func (m Medium) URI() string {
   return m.Raw_URI
}

type Mixed interface {
   Ext() string
   URI() string
}

func (Stream) Ext() string {
   return ".m4v"
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
   c := strings.Join(a, " ")
   if b != "" {
      c += "\n  " + b
   }
   return c
}

func (s Stream) URI() string {
   return s.Raw_URI
}

type Medium struct {
   Characteristics string
   Group_ID string
   Name string
   Raw_URI string
   Type string
}

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   Raw_URI string
}

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

type Streams []Stream

func filter[T Mixed](slice []T, callback func(T) bool) []T {
   var carry []T
   for _, item := range slice {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func index[T Mixed](slice []T, callback func(T, T) bool) int {
   carry := -1
   for i, item := range slice {
      if carry == -1 || callback(slice[carry], item) {
         carry = i
      }
   }
   return carry
}

func (m Media) Filter(f func(Medium) bool) Media {
   return filter(m, f)
}

func (s Streams) Filter(f func(Stream) bool) Streams {
   return filter(s, f)
}

func (m Media) Index(f func(a, b Medium) bool) int {
   return index(m, f)
}

func (s Streams) Index(f func(a, b Stream) bool) int {
   return index(s, f)
}

func (s Streams) Bandwidth(v int) int {
   distance := func(a Stream) int {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return s.Index(func(carry, item Stream) bool {
      return distance(item) < distance(carry)
   })
}
