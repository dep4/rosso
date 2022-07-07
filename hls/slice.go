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

type Item interface {
   Media | Stream
}

type Filter[T Item] func(T) bool

type Master struct {
   Media Slice[Media]
   Stream Slice[Stream]
}

type Media struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
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

type Reduce[T Item] func(T, T) bool

type Slice[T Item] []T

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

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
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
