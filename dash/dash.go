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

func (self Representation) String() string {
   var (
      a []string
      b []string
   )
   if self.Width >= 1 {
      a = append(a, "Width:" + strconv.Itoa(self.Width))
   }
   if self.Height >= 1 {
      a = append(a, "Height:" + strconv.Itoa(self.Height))
   }
   if self.Bandwidth >= 1 {
      a = append(a, "Bandwidth:" + strconv.Itoa(self.Bandwidth))
   }
   b = append(b, "MimeType:" + self.MimeType)
   if self.Codecs != "" {
      b = append(b, "Codecs:" + self.Codecs)
   }
   if self.Adaptation.Lang != "" {
      b = append(b, "Lang:" + self.Adaptation.Lang)
   }
   if self.Adaptation.Role != nil {
      b = append(b, "Role:" + self.Adaptation.Role.Value)
   }
   c := "ID:" + self.ID
   if a != nil {
      c += "\n  " + strings.Join(a, " ")
   }
   return c + "\n  " + strings.Join(b, " ")
}

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

func (self Representation) Ext() string {
   switch self.MimeType {
   case "video/mp4":
      return ".m4v"
   case "audio/mp4":
      return ".m4a"
   }
   return ""
}

func (self Representation) Initialization() string {
   return self.replace_ID(self.SegmentTemplate.Initialization)
}

func (self Representation) Media() []string {
   var (
      media []string
      start int
   )
   if self.SegmentTemplate.StartNumber != nil {
      start = *self.SegmentTemplate.StartNumber
   }
   for _, seg := range self.SegmentTemplate.SegmentTimeline.S {
      for seg.Time = start; seg.Repeat >= 0; seg.Repeat-- {
         medium := self.replace_ID(self.SegmentTemplate.Media)
         time_attr := strconv.Itoa(seg.Time)
         if self.SegmentTemplate.StartNumber != nil {
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

func (self Representation) Role() string {
   if self.Adaptation.Role == nil {
      return ""
   }
   return self.Adaptation.Role.Value
}

func (self Representation) replace_ID(s string) string {
   return strings.Replace(s, "$RepresentationID$", self.ID, 1)
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
