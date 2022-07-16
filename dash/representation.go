package dash

type Representations []Representation

func (self Presentation) Representation() Representations {
   var reps []Representation
   for i, ada := range self.Period.AdaptationSet {
      for _, rep := range ada.Representation {
         rep.Adaptation = &self.Period.AdaptationSet[i]
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

func (self Representations) Filter(callback Filter) Representations {
   var carry []Representation
   for _, item := range self {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (self Representations) Video() Representations {
   return self.Filter(func(a Representation) bool {
      return a.MimeType == "video/mp4"
   })
}

func (self Representations) Audio() Representations {
   return self.Filter(func(a Representation) bool {
      return a.MimeType == "audio/mp4"
   })
}

type Index func(carry, item Representation) bool

func (self Representations) Index(callback Index) int {
   carry := -1
   for i, item := range self {
      if carry == -1 || callback(self[carry], item) {
         carry = i
      }
   }
   return carry
}

func (self Representations) Bandwidth(v int) int {
   distance := func(a Representation) int {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return self.Index(func(carry, item Representation) bool {
      return distance(item) < distance(carry)
   })
}
