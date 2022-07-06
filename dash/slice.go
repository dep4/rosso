package dash

type Representations []Representation

func (m Media) Representations() Representations {
   var reps Representations
   for i, ada := range m.Period.AdaptationSet {
      for _, rep := range ada.Representation {
         rep.Adaptation = &m.Period.AdaptationSet[i]
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

type Filter func(Representation) bool

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

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int `xml:"bandwidth,attr"`
   ContentProtection *ContentProtection
   Height int `xml:"height,attr"`
   SegmentTemplate *SegmentTemplate
   Width int `xml:"width,attr"`
   MimeType string `xml:"mimeType,attr"`
   Codecs string `xml:"codecs,attr"`
   ID string `xml:"id,attr"`
}

func Audio(r Representation) bool {
   return r.MimeType == "audio/mp4"
}

func Video(r Representation) bool {
   return r.MimeType == "video/mp4"
}

func AudioVideo(r Representation) bool {
   return Audio(r) || Video(r)
}
