package dash

import (
   "encoding/xml"
   "io"
   "strconv"
)

// func (a AdaptationSet) Timeline(rep Representation) []string {
func (a AdaptationSet) Timeline() []string {
   var (
      t int
      ts []string
   )
   for _, s := range a.SegmentTemplate.SegmentTimeline.S {
      if s.R == 0 {
         s.R = 1
      }
      for s.R >= 1 {
         ts = append(ts, strconv.Itoa(t))
         t += s.D
         s.R--
      }
   }
   return ts
}

type AdaptationSet struct {
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

func AdaptationSets(src io.Reader) ([]AdaptationSet, error) {
   var mpd struct {
      Period struct {
         AdaptationSet []AdaptationSet
      }
   }
   err := xml.NewDecoder(src).Decode(&mpd)
   if err != nil {
      return nil, err
   }
   return mpd.Period.AdaptationSet, nil
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

func (a AdaptationSet) Main() bool {
   return a.Role.Value == "main"
}
