package dash

import (
   "encoding/xml"
   "io"
   "net/url"
   "strconv"
   "strings"
)

type Adaptation struct {
   MimeType string `xml:"mimeType,attr"`
   Representation []Representation
   SegmentTemplate Template
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

func (r Representation) replace(in string) string {
   return strings.Replace(in, "$RepresentationID$", r.ID, 1)
}

type Segment struct {
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
   T int `xml:"t,attr"`
}

func (s Segment) replace(in string) string {
   return strings.Replace(in, "$Time$", strconv.Itoa(s.T), 1)
}

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
}

func (t Template) Replace(rep Representation) Template {
   t.Initialization = rep.replace(t.Initialization)
   t.Media = rep.replace(t.Media)
   return t
}

func (t Template) URLs(base *url.URL) ([]*url.URL, error) {
   var start int
   addr, err := base.Parse(t.Initialization)
   if err != nil {
      return nil, err
   }
   addrs := []*url.URL{addr}
   for _, seg := range t.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         addr, err := base.Parse(seg.replace(t.Media))
         if err != nil {
            return nil, err
         }
         addrs = append(addrs, addr)
         start += seg.D
         seg.T += seg.D
      }
   }
   return addrs, nil
}
