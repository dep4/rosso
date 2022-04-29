package dash

import (
   "encoding/xml"
   "io"
   "net/url"
   "strconv"
   "strings"
)

func Adaptations(addr *url.URL, body io.Reader) ([]Adaptation, error) {
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

func (a Adaptation) Timeline(rep Representation) []string {
   var (
      meds []string
      t int
   )
   for _, s := range a.SegmentTemplate.SegmentTimeline.S {
      if s.R == 0 {
         s.R = 1
      }
      for s.R >= 1 {
         med := a.SegmentTemplate.Media
         med = strings.Replace(med, "$RepresentationID$", rep.ID, 1)
         med = strings.Replace(med, "$Time$", strconv.Itoa(t), 1)
         meds = append(meds, med)
         t += s.D
         s.R--
      }
   }
   return meds
}

type Adaptation struct {
   SegmentTemplate struct {
      Media string `xml:"media,attr"`
      SegmentTimeline struct {
         S []struct {
            D int `xml:"d,attr"`
            R int `xml:"r,attr"`
         }
      }
   }
   Role struct {
      Value string `xml:"value,attr"`
   }
   Representation []Representation
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
