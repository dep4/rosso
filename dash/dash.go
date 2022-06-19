package dash

import (
   "net/url"
   "strconv"
   "strings"
)

type Protection struct {
   Default_KID string `xml:"default_KID,attr"`
}

type Represent struct {
   ID string `xml:"id,attr"`
   Width int64 `xml:"width,attr"`
   Height int64 `xml:"height,attr"`
   Bandwidth int64 `xml:"bandwidth,attr"` // handle duplicate height
   Codecs string `xml:"codecs,attr"` // handle missing height
   MIME_Type string `xml:"mimeType,attr"`
   ContentProtection *Protection
   SegmentTemplate *Template
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

func (r Represent) Initialization(base *url.URL) (*url.URL, error) {
   ref := r.id(r.SegmentTemplate.Initialization)
   return base.Parse(ref)
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

func (r Represent) id(in string) string {
   return strings.Replace(in, "$RepresentationID$", r.ID, 1)
}

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
   StartNumber *int `xml:"startNumber,attr"`
}

type Segment struct {
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
   T int `xml:"t,attr"`
}

func (s Segment) number(in string) string {
   return strings.Replace(in, "$Number$", strconv.Itoa(s.T), 1)
}

func (s Segment) time(in string) string {
   return strings.Replace(in, "$Time$", strconv.Itoa(s.T), 1)
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
   return r.MIME_Type == "video/mp4"
}

type Adaptation struct {
   ContentProtection *Protection
   Lang string `xml:"lang,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   Representation Represents
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *Template
}

func (m Media) Protection() *Protection {
   for _, ada := range m.Period.AdaptationSet {
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

func Audio(a Adaptation, r Represent) bool {
   if !strings.HasPrefix(a.Lang, "en") {
      return false
   }
   if r.MIME_Type != "audio/mp4" {
      return false
   }
   if a.Role != nil && a.Role.Value != "main" {
      return false
   }
   return true
}

type AdaptationFunc func(Adaptation, Represent) bool

type Media struct {
   Period struct {
      AdaptationSet []Adaptation
   }
}

func (m Media) Represents(fn AdaptationFunc) Represents {
   var reps Represents
   for _, ada := range m.Period.AdaptationSet {
      for _, rep := range ada.Representation {
         if rep.MIME_Type == "" {
            rep.MIME_Type = ada.MIME_Type
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
