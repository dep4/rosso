package dash

type Representations []Representation

type Representation struct {
   Bandwidth int `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Height int `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType string `xml:"mimeType,attr"`
   Width int `xml:"width,attr"`
}

type Filter func(Representation) bool

func (r Representations) Filter(callback Filter) Representations {
   if callback == nil {
      return r
   }
   var carry Representations
   for _, item := range r {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

type Map func(Representation) Representation

func (r Representations) Map(callback Map) Representations {
   if callback == nil {
      return r
   }
   for i, item := range r {
      r[i] = callback(item)
   }
   return r
}

func (r Representations) Reduce(callback Reduce) *Representation {
   if callback == nil {
      return nil
   }
   var carry *Representation
   for _, item := range r {
      carry = callback(carry, item)
   }
   return carry
}

type Reduce func(*Representation, Representation) *Representation
