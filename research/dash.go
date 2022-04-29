package dash

import (
   "encoding/xml"
   "io"
   "net/url"
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
}
