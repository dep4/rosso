package hls

import (
   "sort"
)

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

type Media_Func func(group_ID, name, typ string) int

func (m Media) Media(fn Media_Func) Media {
   if fn == nil {
      return m
   }
   var slice Media
   for _, elem := range m {
      if fn(elem.Group_ID, elem.Name, elem.Type) != 0 {
         slice = append(slice, elem)
      }
   }
   return slice
}

func (m Media) Medium(fn Media_Func) *Medium {
   if len(m) == 0 || fn == nil {
      return nil
   }
   distance := func(i int) int {
      group_ID, name, typ := m[i].Group_ID, m[i].Name, m[i].Type
      return fn(group_ID, name, typ)
   }
   sort.Slice(m, func(a, b int) bool {
      return distance(a) < distance(b)
   })
   return &m[0]
}
