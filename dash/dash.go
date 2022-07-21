package dash

import (
   "strconv"
   "strings"
)

func (r Representation) String() string {
   var (
      a []byte
      space bool
   )
   if r.Width >= 1 {
      a = append(a, "Width:"...)
      a = strconv.AppendInt(a, r.Width, 10)
      a = append(a, " Height:"...)
      a = strconv.AppendInt(a, r.Height, 10)
      space = true
   }
   if r.Bandwidth >= 1 {
      if space {
         a = append(a, ' ')
      }
      a = append(a, "Bandwidth:"...)
      a = strconv.AppendInt(a, r.Bandwidth, 10)
   }
   b := []byte("ID:")
   b = append(b, r.ID...)
   if a != nil {
      b = append(b, "\n  "...)
      b = append(b, a...)
   }
   b = append(b, "\n  MimeType:"...)
   b = append(b, r.MimeType...)
   if r.Codecs != "" {
      b = append(b, " Codecs:"...)
      b = append(b, r.Codecs...)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, " Lang:"...)
      b = append(b, r.Adaptation.Lang...)
   }
   if r.Adaptation.Role != nil {
      b = append(b, " Role:"...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   return string(b)
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
