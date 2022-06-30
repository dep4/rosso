package hls

import (
   "strings"
)

func (m Medium) String() string {
   var b strings.Builder
   b.WriteString("Type:")
   b.WriteString(m.Type)
   b.WriteString(" Name:")
   b.WriteString(m.Name)
   b.WriteString(" ID:")
   b.WriteString(m.Group_ID)
   return b.String()
}

type Medium struct {
   Group_ID string
   Name string
   Raw_URI string
   Type string
}

type Media []Medium

// use GROUP-ID instead
func (m Media) URI(value string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.Raw_URI, value) {
         out = append(out, medium)
      }
   }
   return out
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

func (m Media) Group_ID(value string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.Group_ID, value) {
         out = append(out, medium)
      }
   }
   return out
}

func (m Media) Name(value string) Media {
   var out Media
   for _, medium := range m {
      if medium.Name == value {
         out = append(out, medium)
      }
   }
   return out
}

func (m Media) Type(value string) Media {
   var out Media
   for _, medium := range m {
      if medium.Type == value {
         out = append(out, medium)
      }
   }
   return out
}
