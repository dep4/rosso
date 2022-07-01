package dash

type Adaptation struct {
   Codecs *string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   Lang *string `xml:"lang,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Representation Representations
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
}

type ContentProtection struct {
   Default_KID string `xml:"default_KID,attr"`
}

type Media struct {
   Period struct {
      AdaptationSet []Adaptation
   }
}

func (m Media) Representations() Representations {
   var reps Representations
   for i, ada := range m.Period.AdaptationSet {
      for _, rep := range ada.Representation {
         rep.Adaptation = &m.Period.AdaptationSet[i]
         if rep.Raw_Codecs == nil {
            rep.Raw_Codecs = ada.Codecs
         }
         if rep.ContentProtection == nil {
            rep.ContentProtection = ada.ContentProtection
         }
         if rep.MimeType == nil {
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

type Representations []Representation

type SegmentTemplate struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []struct {
         Duration int `xml:"d,attr"`
         Repeat int `xml:"r,attr"`
         Time int `xml:"t,attr"`
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
}

func (r Representations) Get_Bandwidth(v int64) *Representation {
   distance := func(r *Representation) int64 {
      if r.Bandwidth > v {
         return r.Bandwidth - v
      }
      return v - r.Bandwidth
   }
   var output *Representation
   for i, input := range r {
      if output == nil || distance(&input) < distance(output) {
         output = &r[i]
      }
   }
   return output
}
