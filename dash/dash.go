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
   Representation Representations
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
}

func (r Representation) String() string {
   var (
      a []string
      b []string
      c []string
   )
   a = append(a, "ID:" + r.ID)
   if r.Width >= 1 {
      b = append(b, "Width:" + strconv.Itoa(r.Width))
   }
   if r.Height >= 1 {
      b = append(b, "Height:" + strconv.Itoa(r.Height))
   }
   if r.Bandwidth >= 1 {
      b = append(b, "Bandwidth:" + strconv.Itoa(r.Bandwidth))
   }
   c = append(c, "MimeType:" + r.MimeType)
   if r.Codecs != "" {
      c = append(c, "Codecs:" + r.Codecs)
   }
   if r.Adaptation.Lang != "" {
      c = append(c, "Lang:" + r.Adaptation.Lang)
   }
   if r.Adaptation.Role != nil {
      c = append(c, "Role:" + r.Adaptation.Role.Value)
   }
   s := strings.Join(a, " ")
   if b != nil {
      s += "\n  " + strings.Join(b, " ")
   }
   if c != nil {
      s += "\n  " + strings.Join(c, " ")
   }
   return s
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

type ContentProtection struct {
   Default_KID string `xml:"default_KID,attr"`
}

type Media struct {
   Period struct {
      AdaptationSet []Adaptation
   }
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
