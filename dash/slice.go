package dash

import (
   "strconv"
   "strings"
)

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

type Slice[T Representation] []T

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   Lang string `xml:"lang,attr"`
   MimeType string `xml:"mimeType,attr"`
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Representation []Representation
}

func (s Slice[T]) Filter(callback func(T) bool) Slice[T] {
   var carry []T
   for _, item := range s {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (s Slice[T]) Index(callback func(T, T) bool) int {
   carry := -1
   for i, item := range s {
      if carry == -1 || callback(s[carry], item) {
         carry = i
      }
   }
   return carry
}

func (p Presentation) Representation() Slice[Representation] {
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
