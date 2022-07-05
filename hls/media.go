package hls

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

type Media_Filter func(Medium) bool

func (m Media) Filter(callback Media_Filter) Media {
   var carry Media
   for _, item := range m {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

type Media_Reduce func(*Medium, Medium) *Medium

func (m Media) Reduce(callback Media_Reduce) *Medium {
   var carry *Medium
   for _, item := range m {
      carry = callback(carry, item)
   }
   return carry
}
