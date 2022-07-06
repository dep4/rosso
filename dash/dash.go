package dash

import (
   "strconv"
   "strings"
)

func (r Representation) String() string {
   var (
      s []string
      t []string
   )
   s = append(s, "MimeType:" + r.MimeType)
   if r.Codecs != "" {
      s = append(s, "Codecs:" + r.Codecs)
   }
   if r.Adaptation.Lang != "" {
      s = append(s, "Lang:" + r.Adaptation.Lang)
   }
   if r.Adaptation.Role != nil {
      s = append(s, "Role:" + r.Adaptation.Role.Value)
   }
   if r.Bandwidth >= 1 {
      t = append(t, "Bandwidth:" + strconv.Itoa(r.Bandwidth))
   }
   if r.Width >= 1 {
      t = append(t, "Width:" + strconv.Itoa(r.Width))
   }
   if r.Height >= 1 {
      t = append(t, "Height:" + strconv.Itoa(r.Height))
   }
   js, jt := strings.Join(s, " "), strings.Join(t, " ")
   if jt != "" {
      return js + "\n\t" + jt
   }
   return js
}

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int `xml:"bandwidth,attr"`
   ContentProtection *ContentProtection
   Height int `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType string `xml:"mimeType,attr"`
   Codecs string `xml:"codecs,attr"`
   SegmentTemplate *SegmentTemplate
   Width int `xml:"width,attr"`
}

func (r Representation) Role() string {
   if r.Adaptation.Role == nil {
      return ""
   }
   return r.Adaptation.Role.Value
}

func (r Representation) Ext() string {
   switch r.MimeType {
   case "video/mp4":
      return ".m4v"
   case "audio/mp4":
      return ".m4a"
   case "image/jpeg":
      return ".jpg"
   }
   switch r.Codecs {
   case "stpp":
      return ".ttml"
   case "wvtt":
      return ".vtt"
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

func (r Representation) replace_ID(s string) string {
   return strings.Replace(s, "$RepresentationID$", r.ID, 1)
}
func (r Representations) Bandwidth(v int) *Representation {
   distance := func(r *Representation) int {
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
