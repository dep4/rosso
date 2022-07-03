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
