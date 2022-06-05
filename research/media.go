package hls

import (
   "fmt"
   "net/url"
   "strings"
)

// cdn
func (m Media) RawQuery(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.URI.RawQuery, val) {
         out = append(out, medium)
      }
   }
   return out
}

type Medium struct {
   Type string
   Name string
   GroupID string
   URI *url.URL
}

func (m Media) GroupID(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.GroupID, val) {
         out = append(out, medium)
      }
   }
   return out
}

func (m Medium) Format(f fmt.State, verb rune) {
   fmt.Fprint(f, "Type:", m.Type)
   fmt.Fprint(f, " Name:", m.Name)
   fmt.Fprint(f, " ID:", m.GroupID)
   if verb == 'a' {
      fmt.Fprint(f, " URI:", m.URI)
   }
}

type Media []Medium

// English
func (m Media) Name(val string) Media {
   var out Media
   for _, medium := range m {
      if medium.Name == val {
         out = append(out, medium)
      }
   }
   return out
}

// AUDIO
func (m Media) Type(val string) Media {
   var out Media
   for _, medium := range m {
      if medium.Type == val {
         out = append(out, medium)
      }
   }
   return out
}
