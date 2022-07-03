package hls

import (
   "strings"
)

type Medium struct {
   URI string
   Type string
   Name string
   Group_ID string
   Characteristics string
}

func (m Medium) String() string {
   var b strings.Builder
   b.WriteString("Type:")
   b.WriteString(m.Type)
   b.WriteString(" Name:")
   b.WriteString(m.Name)
   b.WriteString("\n  Group ID:")
   b.WriteString(m.Group_ID)
   if m.Characteristics != "" {
      b.WriteString("\n  Characteristics:")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

func (m Media) Audio() Media {
   var slice Media
   for _, elem := range m {
      if elem.Type == "AUDIO" {
         slice = append(slice, elem)
      }
   }
   return slice
}

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

func (m Media) Ext() string {
   return ".m4a"
}

type Media []Medium
