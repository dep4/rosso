package dash

import (
   "strconv"
   "strings"
)

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs *string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   Height *int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType *string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
   Width *int64 `xml:"width,attr"`
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

func (r Representations) English() *Representation {
   for _, rep := range r {
      if rep.Adaptation.Lang != nil {
         if strings.HasPrefix(*rep.Adaptation.Lang, "en") {
            return &rep
         }
      }
   }
   return nil
}

func (r Representations) AVC1() Representations {
   var reps Representations
   for _, rep := range r {
      if rep.Codecs != nil {
         if strings.HasPrefix(*rep.Codecs, "avc1.") {
            reps = append(reps, rep)
         }
      }
   }
   return reps
}

func (r Representations) MP4A() Representations {
   var reps Representations
   for _, rep := range r {
      if rep.Codecs != nil {
         if strings.HasPrefix(*rep.Codecs, "mp4a.") {
            reps = append(reps, rep)
         }
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
         if rep.Codecs == nil {
            rep.Codecs = ada.Codecs
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

func (r Representation) Ext() string {
   switch *r.MimeType {
   case "audio/mp4":
      return ".m4a"
   case "image/jpeg":
      return ".jpg"
   case "video/mp4":
      return ".m4v"
   }
   switch *r.Codecs {
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
