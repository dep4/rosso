package hls

type Mixed interface {
   Ext() string
   URI() string
}

func (Stream) Ext() string {
   return ".m4v"
}

func (s Stream) URI() string {
   return s.Raw_URI
}

func (Medium) Ext() string {
   return ".m4a"
}

func (m Medium) URI() string {
   return m.Raw_URI
}

type Filter interface {
   Audio(Media) Media
   Audio_Index(Media) int
   Video(Streams) Streams
   Video_Index(Streams) int
}

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
