package dash

import (
   "strconv"
   "strings"
)

func (r Representations) Filter_Codecs(v string) Representations {
   var reps Representations
   for _, rep := range r {
      if rep.Codecs != nil && strings.HasPrefix(*rep.Codecs, v) {
         reps = append(reps, rep)
      }
   }
   return reps
}

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

func (r Representations) Filter_Lang(v string) Representations {
   var reps Representations
   for _, rep := range r {
      var lang string
      if rep.Adaptation.Lang != nil {
         lang = *rep.Adaptation.Lang
      }
      if lang == v {
         reps = append(reps, rep)
      }
   }
   return reps
}

func (r Representations) Filter_Role(v string) Representations {
   var reps Representations
   for _, rep := range r {
      var role string
      if rep.Adaptation.Role != nil {
         role = rep.Adaptation.Role.Value
      }
      if role == v {
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

func (r Representation) String() string {
   var b []byte
   b = append(b, "ID:"...)
   b = append(b, r.ID...)
   if r.Codecs != nil {
      b = append(b, " Codecs:"...)
      b = append(b, *r.Codecs...)
   }
   if r.Adaptation.Lang != nil {
      b = append(b, " Lang:"...)
      b = append(b, *r.Adaptation.Lang...)
   }
   if r.Adaptation.Role != nil {
      b = append(b, " Role:"...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   b = append(b, " Bandwidth:"...)
   b = strconv.AppendInt(b, r.Bandwidth, 10)
   if r.Width != nil {
      b = append(b, " Width:"...)
      b = strconv.AppendInt(b, *r.Width, 10)
      b = append(b, " Height:"...)
      b = strconv.AppendInt(b, *r.Height, 10)
   }
   return string(b)
}

func (r Representation) replace_ID(s string) string {
   return strings.Replace(s, "$RepresentationID$", r.ID, 1)
}

type Representations []Representation

func (r Representations) Reduce_Bandwidth(v int64) *Representation {
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

func (a *Adaptation) Representations() Representations {
   var reps Representations
   for _, rep := range a.Representation {
      rep.Adaptation = a
      if rep.Codecs == nil {
         rep.Codecs = a.Codecs
      }
      if rep.ContentProtection == nil {
         rep.ContentProtection = a.ContentProtection
      }
      if rep.MimeType == nil {
         rep.MimeType = a.MimeType
      }
      if rep.SegmentTemplate == nil {
         rep.SegmentTemplate = a.SegmentTemplate
      }
      reps = append(reps, rep)
   }
   return reps
}

type Representation struct {
   Bandwidth int64 `xml:"bandwidth,attr"`
   ID string `xml:"id,attr"`
   Width *int64 `xml:"width,attr"`
   Height *int64 `xml:"height,attr"`
   Adaptation *Adaptation
   Codecs *string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   MimeType *string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
}

