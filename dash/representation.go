package dash

import (
   "strconv"
   "strings"
)

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   ContentProtection *ContentProtection
   Height *int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Raw_Codecs *string `xml:"codecs,attr"`
   SegmentTemplate *SegmentTemplate
   Width *int64 `xml:"width,attr"`
}

func (r Representation) Role() string {
   if r.Adaptation.Role == nil {
      return ""
   }
   return r.Adaptation.Role.Value
}

func (r Representation) Ext() string {
   switch *r.MimeType {
   case "video/mp4":
      return ".m4v"
   case "audio/mp4":
      return ".m4a"
   case "image/jpeg":
      return ".jpg"
   }
   switch *r.Raw_Codecs {
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
   if r.Raw_Codecs != nil {
      b = append(b, "Codecs:"...)
      b = append(b, *r.Raw_Codecs...)
   } else {
      b = append(b, "MimeType:"...)
      b = append(b, *r.MimeType...)
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
