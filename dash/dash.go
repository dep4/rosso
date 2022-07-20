package dash

import (
   "strconv"
   "strings"
)

func (r Representation) String() string {
   var b strings.Builder
   b.WriteString("ID:")
   b.WriteString(r.ID)
   if line := r.line_two(); line != nil {
      b.WriteByte('\n')
      b.Write(line)
   }
   b.WriteString("\nMimeType:")
   b.WriteString(r.MimeType)
   if r.Codecs != "" {
      b.WriteString(" Codecs:")
      b.WriteString(r.Codecs)
   }
   if r.Adaptation.Lang != "" {
      b.WriteString(" Lang:")
      b.WriteString(r.Adaptation.Lang)
   }
   if r.Adaptation.Role != nil {
      b.WriteString(" Role:")
      b.WriteString(r.Adaptation.Role.Value)
   }
   return b.String()
}

func (r Representation) line_two() []byte {
   var (
      b []byte
      space bool
   )
   if r.Width >= 1 {
      b = append(b, "Width:"...)
      b = strconv.AppendInt(b, r.Width, 10)
      b = append(b, " Height:"...)
      b = strconv.AppendInt(b, r.Height, 10)
      space = true
   }
   if r.Bandwidth >= 1 {
      if space {
         b = append(b, ' ')
      }
      b = append(b, "Bandwidth:"...)
      b = strconv.AppendInt(b, r.Bandwidth, 10)
   }
   return b
}

type Representations []Representation

func (p Presentation) Representation() Representations {
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

type Filter func(Representation) bool

func (r Representations) Filter(callback Filter) Representations {
   var carry []Representation
   for _, item := range r {
      if callback(item) {
         carry = append(carry, item)
      }
   }
   return carry
}

func (r Representations) Video() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MimeType == "video/mp4"
   })
}

func (r Representations) Audio() Representations {
   return r.Filter(func(a Representation) bool {
      return a.MimeType == "audio/mp4"
   })
}

type Index func(carry, item Representation) bool

func (r Representations) Index(callback Index) int {
   carry := -1
   for i, item := range r {
      if carry == -1 || callback(r[carry], item) {
         carry = i
      }
   }
   return carry
}

func (r Representations) Bandwidth(v int64) int {
   distance := func(a Representation) int64 {
      if a.Bandwidth > v {
         return a.Bandwidth - v
      }
      return v - a.Bandwidth
   }
   return r.Index(func(carry, item Representation) bool {
      return distance(item) < distance(carry)
   })
}

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
   Width int64 `xml:"width,attr"`
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
