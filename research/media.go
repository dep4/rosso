package hls

import (
   "fmt"
   "net/url"
)

func (m Media) Query(key, val string) Media {
   var out Media
   for _, medium := range m {
      if medium.URI.Query().Get(key) == val {
         out = append(out, medium)
      }
   }
   return out
}

func (m Media) Name(val string) Media {
   var out Media
   for _, medium := range m {
      if medium.Name == val {
         out = append(out, medium)
      }
   }
   return out
}

func (m Media) Type(val string) Media {
   var out Media
   for _, medium := range m {
      if medium.Type == val {
         out = append(out, medium)
      }
   }
   return out
}

type Media []Medium

type Medium struct {
   Name string
   Type string
   URI *url.URL
}

func (m Medium) Format(f fmt.State, verb rune) {
   fmt.Fprint(f, "Type:", m.Type)
   fmt.Fprint(f, " Name:", m.Name)
   if verb == 'a' {
      fmt.Fprint(f, " URI:", m.URI)
   }
}
