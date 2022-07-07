package dash

import (
   "strconv"
   "strings"
)

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

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   Height int `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
   Width int `xml:"width,attr"`
}

type Representations []Representation

type Filter func(Representation) bool

type Reduce func(Representation, Representation) bool

func (r Representations) Filter(callback Filter) Representations {
   var carry Representations
   for _, item := range r {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (r Representations) Reduce(callback Reduce) *Representation {
   var carry *Representation
   for i, item := range r {
      if carry == nil || callback(*carry, item) {
         carry = &r[i]
      }
   }
   return carry
}

func Audio(r Representation) bool {
   return r.MimeType == "audio/mp4"
}

func Video(r Representation) bool {
   return r.MimeType == "video/mp4"
}

func Bandwidth(v int) func(Representation) int {
   return func(r Representation) int {
      if r.Bandwidth > v {
         return r.Bandwidth - v
      }
      return v - r.Bandwidth
   }
}

func (r Representation) String() string {
   var (
      a []string
      b []string
   )
   if r.Width >= 1 {
      a = append(a, "Width:" + strconv.Itoa(r.Width))
   }
   if r.Height >= 1 {
      a = append(a, "Height:" + strconv.Itoa(r.Height))
   }
   if r.Bandwidth >= 1 {
      a = append(a, "Bandwidth:" + strconv.Itoa(r.Bandwidth))
   }
   b = append(b, "MimeType:" + r.MimeType)
   if r.Codecs != "" {
      b = append(b, "Codecs:" + r.Codecs)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, "Lang:" + r.Adaptation.Lang)
   }
   if r.Adaptation.Role != nil {
      b = append(b, "Role:" + r.Adaptation.Role.Value)
   }
   c := "ID:" + r.ID
   if a != nil {
      c += "\n  " + strings.Join(a, " ")
   }
   return c + "\n  " + strings.Join(b, " ")
}
