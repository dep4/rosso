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

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int `xml:"bandwidth,attr"`
   ContentProtection *ContentProtection
   Height int `xml:"height,attr"`
   SegmentTemplate *SegmentTemplate
   Width int `xml:"width,attr"`
   MimeType string `xml:"mimeType,attr"`
   Codecs string `xml:"codecs,attr"`
   ID string `xml:"id,attr"`
}

func (r Representation) String() string {
   var (
      s []string
      t []string
   )
   if r.Width >= 1 {
      s = append(s, "Width:" + strconv.Itoa(r.Width))
   }
   if r.Height >= 1 {
      s = append(s, "Height:" + strconv.Itoa(r.Height))
   }
   if r.Bandwidth >= 1 {
      s = append(s, "Bandwidth:" + strconv.Itoa(r.Bandwidth))
   }
   if r.Codecs != "" {
      s = append(s, "Codecs:" + r.Codecs)
   }
   t = append(t, "MimeType:" + r.MimeType)
   if r.Adaptation.Lang != "" {
      t = append(t, "Lang:" + r.Adaptation.Lang)
   }
   if r.Adaptation.Role != nil {
      t = append(t, "Role:" + r.Adaptation.Role.Value)
   }
   t = append(t, "ID:" + r.ID)
   js, jt := strings.Join(s, " "), strings.Join(t, " ")
   if jt != "" {
      return js + "\n  " + jt
   }
   return js
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
