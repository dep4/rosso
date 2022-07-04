package hls

import (
   "strings"
)

func (m Media) Get_Group_ID(value string) *Medium {
   for _, medium := range m {
      if medium.Group_ID == value {
         return &medium
      }
   }
   return nil
}

func (m Media) Get_Name(value string) *Medium {
   for _, medium := range m {
      if medium.Name == value {
         return &medium
      }
   }
   return nil
}

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

type Media_Filter interface {
   Group_ID() string
   Name() string
   Type() string
}

func (m Media) Filter(f Media_Filter) Media {
   if f == nil {
      return m
   }
   pass := func(m Medium) bool {
      if !strings.Contains(m.Group_ID, f.Group_ID()) {
         return false
      }
      if f.Name() != "" && f.Name() != m.Name {
         return false
      }
      if f.Type() != "" && f.Type() != m.Type {
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
