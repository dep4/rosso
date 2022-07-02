package dash

import (
   "strings"
)

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   Lang string `xml:"lang,attr"`
   MimeType string `xml:"mimeType,attr"`
   Representation Representations
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
}

func (r Representations) Audio() Representations {
   var reps Representations
   for _, rep := range r {
      if !strings.HasPrefix(rep.Adaptation.Lang, "en") {
         continue
      }
      if rep.MimeType != "audio/mp4" {
         continue
      }
      if rep.Role() == "description" {
         continue
      }
      reps = append(reps, rep)
   }
   return reps
}

func (r Representations) Video() Representations {
   var reps Representations
   for _, rep := range r {
      if rep.MimeType == "video/mp4" {
         reps = append(reps, rep)
      }
   }
   return reps
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
