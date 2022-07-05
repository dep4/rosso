package hls

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

type Media_Filter interface {
   Group_ID(string) bool
   Name(string) bool
   Type(string) bool
}

func (m Media) Media(f Media_Filter) Media {
   var prev Media
   for _, curr := range m {
      if !f.Group_ID(curr.Group_ID) {
         continue
      }
      if !f.Name(curr.Name) {
         continue
      }
      if !f.Type(curr.Type) {
         continue
      }
      prev = append(prev, curr)
   }
   return prev
}

func (m Media) Medium(f Media_Filter) *Medium {
   for _, curr := range m {
      if f.Group_ID(curr.Group_ID) {
         return &curr
      }
      if f.Name(curr.Name) {
         return &curr
      }
      if f.Type(curr.Type) {
         return &curr
      }
   }
   return nil
}
