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
   Role struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate struct {
      Initialization string `xml:"initialization,attr"`
      Media string `xml:"media,attr"`
      SegmentTimeline struct {
         S []struct {
            D int `xml:"d,attr"`
            R int `xml:"r,attr"`
         }
      }
   }
}

func Adaptations(body io.Reader) ([]Adaptation, error) {
   var mpd struct {
      Period struct {
         AdaptationSet []Adaptation
      }
   }
   err := xml.NewDecoder(body).Decode(&mpd)
   if err != nil {
      return nil, err
   }
   return mpd.Period.AdaptationSet, nil
}

func (a Adaptation) URLs(u *url.URL, r Representation) ([]*url.URL, error) {
   var start int
   parse := func(s string, t int) (*url.URL, error) {
      s = strings.Replace(s, "$RepresentationID$", r.ID, 1)
      s = strings.Replace(s, "$Time$", strconv.Itoa(t), 1)
      return u.Parse(s)
   }
   addr, err := parse(a.SegmentTemplate.Initialization, 0)
   if err != nil {
      return nil, err
   }
   addrs := []*url.URL{addr}
   for _, seg := range a.SegmentTemplate.SegmentTimeline.S {
      for seg.R >= 0 {
         addr, err := parse(a.SegmentTemplate.Media, start)
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
