package hls

func (m Media) Filter(callback Media_Filter) Media {
   var carry Media
   for _, item := range m {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (m Media) Reduce(callback Media_Reduce) *Medium {
   var carry *Medium
   for _, item := range m {
      carry = callback(carry, item)
   }
   return carry
}

type Media_Filter func(Medium) bool

type Media_Reduce func(*Medium, Medium) *Medium

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}
