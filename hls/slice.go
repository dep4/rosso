package hls

import (
   "strconv"
   "strings"
)

func Bandwidth(v int) func(Stream) int {
   return func(s Stream) int {
      if s.Bandwidth > v {
         return s.Bandwidth - v
      }
      return v - s.Bandwidth
   }
}

type Media struct {
   Characteristics string
   Group_ID string
   Name string
   Raw_URI string
   Type string
}

func (Media) Ext() string {
   return ".m4a"
}

func (m Media) String() string {
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

func (m Media) URI() string {
   return m.Raw_URI
}

type Mixed interface {
   Ext() string
   URI() string
}

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   Raw_URI string
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

type Slice[T Mixed] []T

type Master struct {
   Media Slice[Media]
   Stream Slice[Stream]
}

func (s Slice[T]) Filter(callback func(T) bool) Slice[T] {
   if callback == nil {
      return s
   }
   var carry []T
   for _, item := range s {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (s Slice[T]) Index(callback func(T, T) bool) int {
   carry := -1
   for i, item := range s {
      if carry == -1 || callback(s[carry], item) {
         carry = i
      }
   }
   return carry
}
