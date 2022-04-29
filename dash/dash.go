package dash

import (
   "encoding/xml"
   "io"
   "net/url"
   "strconv"
   "strings"
)

type Segment struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []struct {
         D int `xml:"d,attr"`
         R int `xml:"r,attr"`
      }
   }
}

func replace(s string, r Representation, t int) string {
   s = strings.Replace(s, "$RepresentationID$", r.ID, 1)
   s = strings.Replace(s, "$Time$", strconv.Itoa(t), 1)
   return s
}

func (s Segment) URLs(u *url.URL, r Representation) ([]*url.URL, error) {
   addr, err := u.Parse(replace(s.Initialization, r, 0))
   if err != nil {
      return nil, err
   }
   addrs := []*url.URL{addr}
   var start int
   for _, seg := range s.SegmentTimeline.S {
      for seg.R >= 0 {
         addr, err := u.Parse(replace(s.Media, r, start))
         if err != nil {
            return nil, err
         }
         addrs = append(addrs, addr)
         start += seg.D
         seg.R--
      }
   }
   return addrs, nil
}

type Adaptation struct {
   MimeType string `xml:"mimeType,attr"`
   Representation []Representation
   Role struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate Segment
}

func Adaptations(body io.Reader) ([]Adaptation, error) {
   var media struct {
      Period struct {
         AdaptationSet []Adaptation
      }
   }
   err := xml.NewDecoder(body).Decode(&media)
   if err != nil {
      return nil, err
   }
   return media.Period.AdaptationSet, nil
}

type Representation struct {
   ID string `xml:"id,attr"`
   Width int64 `xml:"width,attr"`
   Height int64 `xml:"height,attr"`
   Bandwidth int64 `xml:"bandwidth,attr"`
}

func (r Representation) String() string {
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
   return string(buf)
}
