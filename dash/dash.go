package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "net/url"
   "strings"
)

type Period struct {
   AdaptationSet []struct {
      ContentType string `xml:"contentType,attr"`
      MimeType string `xml:"mimeType,attr"`
      Representation []Represent
      Role struct {
         Value string `xml:"value,attr"`
      }
      SegmentTemplate Template
   }
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

func (p Period) Audio(bandwidth int) *Represent {
   return p.represent(bandwidth, "audio/mp4")
}

func (p Period) Video(bandwidth int) *Represent {
   return p.represent(bandwidth, "video/mp4")
}

func (p Period) represent(bandwidth int, typ string) *Represent {
   distance := func(r *Represent) int {
      if r.Bandwidth > bandwidth {
         return r.Bandwidth - bandwidth
      }
      return bandwidth - r.Bandwidth
   }
   var dst *Represent
   for i, ada := range p.AdaptationSet {
      if ada.MimeType == typ {
         for j, src := range ada.Representation {
            if j == 0 || distance(&src) < distance(dst) {
               dst = &p.AdaptationSet[i].Representation[j]
            }
         }
      }
   }
   return dst
}

type Represent struct {
   ID string `xml:"id,attr"`
   Width int `xml:"width,attr"`
   Height int `xml:"height,attr"`
   Bandwidth int `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   ContentProtection []struct {
      SchemeIdUri string `xml:"schemeIdUri,attr"`
      PSSH string `xml:"pssh"`
   }
}

func (r Represent) Format(f fmt.State, verb rune) {
   fmt.Fprint(f, "ID:", r.ID)
   if r.Width >= 1 {
      fmt.Fprint(f, " Width:", r.Width)
      fmt.Fprint(f, " Height:", r.Height)
   }
   fmt.Fprint(f, " Bandwidth:", r.Bandwidth)
   fmt.Fprint(f, " Codec:", r.Codecs)
   if verb == 'a' {
      for _, con := range r.ContentProtection {
         fmt.Fprint(f, "\nSchemeID:", con.SchemeIdUri)
         if con.PSSH != "" {
            fmt.Fprint(f, "\nPSSH:", con.PSSH)
         }
      }
   }
}

func (r Represent) replace(in string) string {
   return strings.Replace(in, "$RepresentationID$", r.ID, 1)
}

type Segment struct {
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
   T int `xml:"t,attr"`
}

func (s Segment) replace(in string) string {
   return strings.Replace(in, "$Time$", fmt.Sprint(s.T), 1)
}

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
}

func (t Template) Base(rep Represent) string {
   return rep.replace(t.Initialization)
}

func (t Template) URL(rep Represent, base *url.URL) ([]*url.URL, error) {
   var start int
   addr, err := base.Parse(rep.replace(t.Initialization))
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
