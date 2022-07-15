package dash

type Representations []Representation

func (p Presentation) Representation() Representations {
   var reps []Representation
   for i, ada := range p.Period.AdaptationSet {
      for _, rep := range ada.Representation {
         rep.Adaptation = &p.Period.AdaptationSet[i]
         if rep.Codecs == "" {
            rep.Codecs = ada.Codecs
         }
         if rep.ContentProtection == nil {
            rep.ContentProtection = ada.ContentProtection
         }
         if rep.MimeType == "" {
            rep.MimeType = ada.MimeType
         }
         if rep.SegmentTemplate == nil {
            rep.SegmentTemplate = ada.SegmentTemplate
         }
         reps = append(reps, rep)
      }
   }
   return reps
}

type Filter func(Representation) bool

func (r Representations) Filter(callback Filter) Representations {
   var carry []Representation
   for _, item := range r {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (r Representations) Video() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MimeType == "video/mp4"
   })
}

func (r Representations) Audio() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MimeType == "audio/mp4"
   })
}

type Index func(carry, item Representation) bool

func (r Representations) Index(callback Index) int {
   carry := -1
   for i, item := range r {
      if carry == -1 || callback(r[carry], item) {
         carry = i
      }
   }
   return carry
}

func (r Representations) Bandwidth(v int) int {
   distance := func(a Representation) int {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return r.Index(func(carry, item Representation) bool {
      return distance(item) < distance(carry)
   })
}
