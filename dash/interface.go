package dash

func Audio(r []Representation) []Representation {
   var carry []Representation
   for _, item := range r {
      if item.MimeType == "audio/mp4" {
         carry = append(carry, item)
      }
   }
   return carry
}

func Bandwidth(r []Representation, b int) int {
   distance := func(r Representation) int {
      if r.Bandwidth > b {
         return r.Bandwidth - b
      }
      return b - r.Bandwidth
   }
   carry := -1
   for i, item := range r {
      if carry == -1 || distance(item) < distance(r[carry]) {
         carry = i
      }
   }
   return carry
}

func Video(r []Representation) []Representation {
   var carry []Representation
   for _, item := range r {
      if item.MimeType == "video/mp4" {
         carry = append(carry, item)
      }
   }
   return carry
}

type Filter interface {
   Audio([]Representation) []Representation
   Audio_Index([]Representation) int
   Video([]Representation) []Representation
   Video_Index([]Representation) int
}
