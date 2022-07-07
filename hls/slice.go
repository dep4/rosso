package hls

import (
   "strconv"
   "strings"
)

type Element interface {
   Media | Stream
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
   ja := strings.Join(a, " ")
   if b != "" {
      return ja + "\n  " + b
   }
   return ja
}

type Filter[T Element] func(T) bool

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

type Reduce[T Element] func(*T, T) *T

func Bandwidth(v int) Reduce[Stream] {
   distance := func(s *Stream) int {
      if s.Bandwidth > v {
         return s.Bandwidth - v
      }
      return v - s.Bandwidth
   }
   return func(carry *Stream, item Stream) *Stream {
      if carry == nil || distance(&item) < distance(carry) {
         return &item
      }
      return carry
   }
}

func Name(v string) Reduce[Media] {
   return func(carry *Media, item Media) *Media {
      if item.Name == v {
         return &item
      }
      return carry
   }
}

type Slice[T Element] []T

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
   if callback == nil {
      return nil
   }
   var carry *T
   for _, item := range s {
      carry = callback(carry, item)
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
