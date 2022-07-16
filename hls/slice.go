package hls

import (
   "strconv"
   "strings"
)

func (Medium) Ext() string {
   return ".m4a"
}

func (self Medium) String() string {
   var b strings.Builder
   b.WriteString("Type:")
   b.WriteString(self.Type)
   b.WriteString(" Name:")
   b.WriteString(self.Name)
   b.WriteString("\n  Group ID:")
   b.WriteString(self.Group_ID)
   if self.Characteristics != "" {
      b.WriteString("\n  Characteristics:")
      b.WriteString(self.Characteristics)
   }
   return b.String()
}

func (self Medium) URI() string {
   return self.Raw_URI
}

type Mixed interface {
   Ext() string
   URI() string
}

func (Stream) Ext() string {
   return ".m4v"
}

func (self Stream) String() string {
   var (
      a []string
      b string
   )
   if self.Resolution != "" {
      a = append(a, "Resolution:" + self.Resolution)
   }
   a = append(a, "Bandwidth:" + strconv.Itoa(self.Bandwidth))
   if self.Codecs != "" {
      a = append(a, "Codecs:" + self.Codecs)
   }
   if self.Audio != "" {
      b = "Audio:" + self.Audio
   }
   c := strings.Join(a, " ")
   if b != "" {
      c += "\n  " + b
   }
   return c
}

func (self Stream) URI() string {
   return self.Raw_URI
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

func (self Media) Filter(f func(Medium) bool) Media {
   return filter(self, f)
}

func (self Streams) Filter(f func(Stream) bool) Streams {
   return filter(self, f)
}

func (self Media) Index(f func(a, b Medium) bool) int {
   return index(self, f)
}

func (self Streams) Index(f func(a, b Stream) bool) int {
   return index(self, f)
}

func (self Streams) Bandwidth(v int) int {
   distance := func(a Stream) int {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return self.Index(func(carry, item Stream) bool {
      return distance(item) < distance(carry)
   })
}
