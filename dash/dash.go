package dash

import (
   "encoding/hex"
   "encoding/xml"
   "io"
   "net/url"
   "strconv"
   "strings"
)

func (r Represent) Media(base *url.URL) ([]*url.URL, error) {
   var addrs []*url.URL
   start := r.SegmentTemplate.StartNumber
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := r.id(r.SegmentTemplate.Media)
         if start >= 1 {
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

const (
   Audio = "audio/mp4"
   Video = "video/mp4"
)

type Protection struct {
   Default_KID string `xml:"default_KID,attr"`
}

func (p Protection) KID() ([]byte, error) {
   kid := strings.ReplaceAll(p.Default_KID, "-", "")
   return hex.DecodeString(kid)
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

type Represents []Represent

func (r Represents) Represent(bandwidth int64) *Represent {
   distance := func(r *Represent) int64 {
      if r.Bandwidth > bandwidth {
         return r.Bandwidth - bandwidth
      }
      return bandwidth - r.Bandwidth
   }
   var dst *Represent
   for i, src := range r {
      if dst == nil || distance(&src) < distance(dst) {
         dst = &r[i]
      }
   }
   return dst
}

type Period struct {
   AdaptationSet []struct {
      ContentProtection *Protection
      MimeType string `xml:"mimeType,attr"`
      Representation []Represent
      Role *struct {
         Value string `xml:"value,attr"`
      }
      SegmentTemplate *Template
   }
}

func (p Period) MimeType(typ string) Represents {
   var reps Represents
   for _, ada := range p.AdaptationSet {
      if ada.Role == nil {
         for _, rep := range ada.Representation {
            if ada.MimeType == typ || rep.MimeType == typ {
               if rep.SegmentTemplate == nil {
                  rep.SegmentTemplate = ada.SegmentTemplate
               }
               reps = append(reps, rep)
            }
         }
      }
   }
   return reps
}

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
   StartNumber int `xml:"startNumber,attr"`
}
