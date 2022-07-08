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

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   Raw_URI string
}

type Media struct {
   Raw_URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

type Item interface {
   Ext() string
   URI() string
}

type Filter[T Item] func(T) bool

type Reduce[T Item] func(T, T) bool

func (m Media) URI() string {
   return m.Raw_URI
}

func (s Stream) URI() string {
   return s.Raw_URI
}

func (Media) Ext() string {
   return ".m4a"
}

func (Stream) Ext() string {
   return ".m4v"
}

func (s Slice[T]) Filter(callback Filter[T]) Slice[T] {
   if callback == nil {
      return s
   }
   var carry Slice[T]
   for _, item := range s {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (s Slice[T]) Reduce(callback Reduce[T]) *T {
   var carry *T
   for i, item := range s {
      if carry == nil || callback(*carry, item) {
         carry = &s[i]
      }
   }
   return carry
}

type Master struct {
   Media Slice[Media]
   Stream Slice[Stream]
}

type Slice[T Item] []T
