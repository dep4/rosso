package dash

type Filter interface {
   Audio(Representations) Representations
   Audio_Index(Representations) int
   Video(Representations) Representations
   Video_Index(Representations) int
}

func (r Representations) Filter(f func(Representation) bool) Representations {
   var carry []Representation
   for _, item := range r {
      if f(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (r Representations) Video() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MimeType == "video/mp4"
   })
}

func (r Representations) Audio() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MimeType == "audio/mp4"
   })
}

func (r Representations) Index(f func(a, b Representation) bool) int {
   carry := -1
   for i, item := range r {
      if carry == -1 || f(r[carry], item) {
         carry = i
      }
   }
   return carry
}

func (r Representations) Bandwidth(v int) int {
   distance := func(a Representation) int {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return r.Index(func(carry, item Representation) bool {
      return distance(item) < distance(carry)
   })
}
