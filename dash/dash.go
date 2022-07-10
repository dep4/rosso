package dash

import (
   "strconv"
   "strings"
)

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

type ContentProtection struct {
   Default_KID string `xml:"default_KID,attr"`
}

type Presentation struct {
   Period struct {
      AdaptationSet []Adaptation
   }
}

func (p Presentation) Representation() []Representation {
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

func (r Representation) Ext() string {
   switch r.MimeType {
   case "video/mp4":
      return ".m4v"
   case "audio/mp4":
      return ".m4a"
   }
   return ""
}

func (r Representation) Initialization() string {
   return r.replace_ID(r.SegmentTemplate.Initialization)
}

func (r Representation) Media() []string {
   var (
      media []string
      start int
   )
   if r.SegmentTemplate.StartNumber != nil {
      start = *r.SegmentTemplate.StartNumber
   }
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.Time = start; seg.Repeat >= 0; seg.Repeat-- {
         medium := r.replace_ID(r.SegmentTemplate.Media)
         time_attr := strconv.Itoa(seg.Time)
         if r.SegmentTemplate.StartNumber != nil {
            medium = strings.Replace(medium, "$Number$", time_attr, 1)
            seg.Time++
            start++
         } else {
            medium = strings.Replace(medium, "$Time$", time_attr, 1)
            seg.Time += seg.Duration
            start += seg.Duration
         }
         media = append(media, medium)
      }
   }
   return media
}

func (r Representation) Role() string {
   if r.Adaptation.Role == nil {
      return ""
   }
   return r.Adaptation.Role.Value
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

func (r Representation) replace_ID(s string) string {
   return strings.Replace(s, "$RepresentationID$", r.ID, 1)
}

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
