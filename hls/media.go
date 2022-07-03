package hls

import (
   "bytes"
   "strings"
)

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

func (m Media) Ext(b []byte) string {
   if bytes.Contains(b, []byte("ftypiso5")) {
      return ".m4a"
   }
   if bytes.HasPrefix(b, []byte{'G'}) {
      return ".mts"
   }
   return ""
}

type Medium struct {
   Group_ID string
   Name string
   Type string
   URI string
}

func (m Medium) String() string {
   var b strings.Builder
   b.WriteString("Type:")
   b.WriteString(m.Type)
   b.WriteString(" Name:")
   b.WriteString(m.Name)
   b.WriteString("\n\tGROUP-ID:")
   b.WriteString(m.Group_ID)
   return b.String()
}

type Media []Medium
