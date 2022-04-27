package dash

import (
   "encoding/xml"
   "io"
   "strconv"
)

type AdaptationSet struct {
   Role struct {
      Value string `xml:"value,attr"`
   }
   Representation []Representation
}

type Media struct {
   Period struct {
      AdaptationSet []AdaptationSet
   }
}

func NewMedia(src io.Reader) (*Media, error) {
   med := new(Media)
   err := xml.NewDecoder(src).Decode(med)
   if err != nil {
      return nil, err
   }
   return med, nil
}

func (m Media) Main() []AdaptationSet {
   var adas []AdaptationSet
   for _, ada := range m.Period.AdaptationSet {
      if ada.Role.Value == "main" {
         adas = append(adas, ada)
      }
   }
   return adas
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
