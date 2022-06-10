package dash

import (
   "encoding/hex"
   "encoding/xml"
   "io"
   "net/url"
   "strconv"
   "strings"
)

func Audio(a Adaptation, r Represent) bool {
   if !strings.HasPrefix(a.Lang, "en") {
      return false
   }
   if r.MimeType != "audio/mp4" {
      return false
   }
   if a.Role != nil && a.Role.Value != "main" {
      return false
   }
   return true
}

type Adaptation struct {
   ContentProtection *Protection
   Lang string `xml:"lang,attr"`
   MimeType string `xml:"mimeType,attr"`
   Representation []Represent
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *Template
}

type Represent struct {
   ID string `xml:"id,attr"` // RepresentationID
   Width int64 `xml:"width,attr"`
   Height int64 `xml:"height,attr"`
   Bandwidth int64 `xml:"bandwidth,attr"` // handle duplicate height
   Codecs string `xml:"codecs,attr"` // handle missing height
   MimeType string `xml:"mimeType,attr"`
   ContentProtection *Protection
   SegmentTemplate *Template
}

func (p Period) Represents(fn PeriodFunc) Represents {
   var reps Represents
   for _, ada := range p.AdaptationSet {
      for _, rep := range ada.Representation {
         if rep.MimeType == "" {
            rep.MimeType = ada.MimeType
         }
         if rep.SegmentTemplate == nil {
            rep.SegmentTemplate = ada.SegmentTemplate
         }
         if fn(ada, rep) {
            reps = append(reps, rep)
         }
      }
   }
   return reps
}

type Period struct {
   AdaptationSet []Adaptation
}

type PeriodFunc func(Adaptation, Represent) bool

type Protection struct {
   Default_KID string `xml:"default_KID,attr"`
}

func (p Protection) KeyID() ([]byte, error) {
   keyID := strings.ReplaceAll(p.Default_KID, "-", "")
   return hex.DecodeString(keyID)
}

func (r Represent) String() string {
   var buf []byte
   buf = append(buf, "ID:"...)
   buf = append(buf, r.ID...)
   if r.Width >= 1 {
      buf = append(buf, " Width:"...)
      buf = strconv.AppendInt(buf, r.Width, 10)
      buf = append(buf, " Height:"...)
      buf = strconv.AppendInt(buf, r.Height, 10)
   }
   buf = append(buf, " Bandwidth:"...)
   buf = strconv.AppendInt(buf, r.Bandwidth, 10)
   buf = append(buf, " Codec:"...)
   buf = append(buf, r.Codecs...)
   return string(buf)
}

type Segment struct {
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
   T int `xml:"t,attr"`
}

func (r Represent) Initialization(base *url.URL) (*url.URL, error) {
   ref := r.id(r.SegmentTemplate.Initialization)
   return base.Parse(ref)
}

func NewPeriod(body io.Reader) (*Period, error) {
   var media struct {
      Period Period
   }
   err := xml.NewDecoder(body).Decode(&media)
   if err != nil {
      return nil, err
   }
   return &media.Period, nil
}

func (p Period) Protection() *Protection {
   for _, ada := range p.AdaptationSet {
      if ada.ContentProtection != nil {
         return ada.ContentProtection
      }
      for _, rep := range ada.Representation {
         if rep.ContentProtection != nil {
            return rep.ContentProtection
         }
      }
   }
   return nil
}

type Represents []Represent

func (r Represents) Represent(bandwidth int64) *Represent {
   distance := func(r *Represent) int64 {
      if r.Bandwidth > bandwidth {
         return r.Bandwidth - bandwidth
      }
      return bandwidth - r.Bandwidth
   }
   var output *Represent
   for i, input := range r {
      if output == nil || distance(&input) < distance(output) {
         output = &r[i]
      }
   }
   return output
}

func Video(a Adaptation, r Represent) bool {
   return r.MimeType == "video/mp4"
}

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
   StartNumber *int `xml:"startNumber,attr"`
}

func (r Represent) Media(base *url.URL) ([]*url.URL, error) {
   var (
      addrs []*url.URL
      start int
   )
   if r.SegmentTemplate.StartNumber != nil {
      start = *r.SegmentTemplate.StartNumber
   }
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := r.id(r.SegmentTemplate.Media)
         if r.SegmentTemplate.StartNumber != nil {
            ref = seg.number(ref)
            seg.T++
            start++
         } else {
            ref = seg.time(ref)
            seg.T += seg.D
            start += seg.D
         }
         addr, err := base.Parse(ref)
         if err != nil {
            return nil, err
         }
         addrs = append(addrs, addr)
      }
   }
   return addrs, nil
}
