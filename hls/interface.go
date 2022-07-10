package hls

func Bandwidth(s []Stream, b int) int {
   distance := func(s Stream) int {
      if s.Bandwidth > b {
         return s.Bandwidth - b
      }
      return b - s.Bandwidth
   }
   carry := -1
   for i, item := range s {
      if carry == -1 || distance(item) < distance(s[carry]) {
         carry = i
      }
   }
   return carry
}

type Filter interface {
   Audio([]Media) []Media
   Audio_Index([]Media) int
   Video([]Stream) []Stream
   Video_Index([]Stream) int
}

func (Media) Ext() string {
   return ".m4a"
}

func (m Media) URI() string {
   return m.Raw_URI
}

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
