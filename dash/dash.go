package dash

import (
   "encoding/hex"
   "encoding/xml"
   "io"
   "net/url"
   "strconv"
   "strings"
)

const (
   Audio = "audio/mp4"
   Video = "video/mp4"
)

type Adaptation struct {
   MimeType string `xml:"mimeType,attr"`
   Representation []Represent
   ContentProtection *Protection
   SegmentTemplate *Template
}

type AdaptationSet []Adaptation

func NewAdaptationSet(body io.Reader) (AdaptationSet, error) {
   var media struct {
      Period struct {
         AdaptationSet AdaptationSet
      }
   }
   err := xml.NewDecoder(body).Decode(&media)
   if err != nil {
      return nil, err
   }
   return media.Period.AdaptationSet, nil
}

func (a AdaptationSet) MimeType(typ string) AdaptationSet {
   var adapts AdaptationSet
   for _, adapt := range a {
      if adapt.MimeType == typ {
         adapts = append(adapts, adapt)
      }
   }
   return adapts
}

func (a AdaptationSet) Protection() *Protection {
   for _, adapt := range a {
      if adapt.ContentProtection != nil {
         return adapt.ContentProtection
      }
      for _, rep := range adapt.Representation {
         if rep.ContentProtection != nil {
            return rep.ContentProtection
         }
      }
   }
   return nil
}

func (a AdaptationSet) Represent(bandwidth int64) *Represent {
   distance := func(r *Represent) int64 {
      if r.Bandwidth > bandwidth {
         return r.Bandwidth - bandwidth
      }
      return bandwidth - r.Bandwidth
   }
   var dst *Represent
   for i, adapt := range a {
      for j, src := range adapt.Representation {
         if dst == nil || distance(&src) < distance(dst) {
            dst = &a[i].Representation[j]
            if dst.SegmentTemplate == nil {
               dst.SegmentTemplate = adapt.SegmentTemplate
            }
         }
      }
   }
   return dst
}

type Protection struct {
   Default_KID string `xml:"default_KID,attr"`
}

func (p Protection) KID() ([]byte, error) {
   kid := strings.ReplaceAll(p.Default_KID, "-", "")
   return hex.DecodeString(kid)
}

type Represent struct {
   ID string `xml:"id,attr"` // RepresentationID
   Width int64 `xml:"width,attr"`
   Height int64 `xml:"height,attr"`
   Bandwidth int64 `xml:"bandwidth,attr"` // handle duplicate height
   Codecs string `xml:"codecs,attr"` // handle missing height
   ContentProtection *Protection
   SegmentTemplate *Template
}

func (r Represent) Initialization(base *url.URL) (*url.URL, error) {
   ref := r.id(r.SegmentTemplate.Initialization)
   return base.Parse(ref)
}

func (r Represent) Media(base *url.URL) ([]*url.URL, error) {
   var addrs []*url.URL
   start, number := r.number()
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := r.id(r.SegmentTemplate.Media)
         if number {
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

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
   StartNumber *int `xml:"startNumber,attr"`
}
