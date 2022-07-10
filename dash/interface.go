package dash

func Audio(r Representations) Representations {
   var carry Representations
   for _, item := range r {
      if item.MimeType == "audio/mp4" {
         carry = append(carry, item)
      }
   }
   return carry
}

func Bandwidth(r Representations, b int) int {
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

func Video(r Representations) Representations {
   var carry Representations
   for _, item := range r {
      if item.MimeType == "video/mp4" {
         carry = append(carry, item)
      }
   }
   return carry
}

type Filter interface {
   Audio(Representations) Representations
   Audio_Index(Representations) int
   Video(Representations) Representations
   Video_Index(Representations) int
}
