package hls

import (
   "sort"
)

func (m Media) Reduce(r Media_Reduce) *Medium {
   if len(m) == 0 {
      return nil
   }
   distance := func(i int) int {
      return r.distance(m[i].Group_ID, m[i].Name)
   }
   sort.Slice(m, func(a, b int) bool {
      return distance(a) < distance(b)
   })
   return &m[0]
}

type Media_Reduce interface {
   distance(group_ID, name string) int
}

type Media_Filter interface {
   Group_ID(string) bool
   Name(string) bool
   Type(string) bool
}

func (m Media) Filter(f Media_Filter) Media {
   if f == nil {
      return m
   }
   pass := func(m Medium) bool {
      if !f.Group_ID(m.Group_ID) {
         return false
      }
      if !f.Name(m.Name) {
         return false
      }
      if !f.Type(m.Type) {
         return false
      }
      return true
   }
   var slice Media
   for _, elem := range m {
      if pass(elem) {
         slice = append(slice, elem)
      }
   }
   return slice
}

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}
