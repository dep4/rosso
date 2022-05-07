package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "net/url"
   "strings"
)

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

func (r Represent) Format(f fmt.State, verb rune) {
   fmt.Fprint(f, "ID:", r.ID)
   if r.Width >= 1 {
      fmt.Fprint(f, " Width:", r.Width)
      fmt.Fprint(f, " Height:", r.Height)
   }
   fmt.Fprint(f, " Bandwidth:", r.Bandwidth)
   fmt.Fprint(f, " Codec:", r.Codecs)
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
            if dst == nil || distance(&src) < distance(dst) {
               dst = &p.AdaptationSet[i].Representation[j]
               if dst.SegmentTemplate == nil {
                  dst.SegmentTemplate = ada.SegmentTemplate
               }
            }
         }
      }
   }
   return dst
}

type Segment struct {
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
   T int `xml:"t,attr"`
}

type Period struct {
   AdaptationSet []struct {
      MimeType string `xml:"mimeType,attr"`
      Representation []Represent
      SegmentTemplate *Template
   }
}

type Represent struct {
   ID string `xml:"id,attr"`
   Width int `xml:"width,attr"`
   Height int `xml:"height,attr"`
   Bandwidth int `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   SegmentTemplate *Template
}

func (r Represent) replace(in string) string {
   return strings.Replace(in, "$RepresentationID$", r.ID, 1)
}

func (r Represent) Time(base *url.URL) ([]*url.URL, error) {
   parse := func(addr string) (*url.URL, error) {
      ref := r.replace(addr)
      return base.Parse(ref)
   }
   addr, err := parse(r.SegmentTemplate.Initialization)
   if err != nil {
      return nil, err
   }
   addrs := []*url.URL{addr}
   var start int
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := seg.replace(r.SegmentTemplate.Media)
         addr, err := parse(ref)
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

func (s Segment) replace(in string) string {
   return strings.Replace(in, "$Time$", fmt.Sprint(s.T), 1)
}

type Template struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []Segment
   }
   StartNumber int `xml:"startNumber,attr"`
}

func (r Represent) Number(base *url.URL) ([]*url.URL, error) {
   parse := func(addr string) (*url.URL, error) {
      ref := r.replace(addr)
      return base.Parse(ref)
   }
   addr, err := parse(r.SegmentTemplate.Initialization)
   if err != nil {
      return nil, err
   }
   addrs := []*url.URL{addr}
   // FIXME
   var start int
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := seg.replace(r.SegmentTemplate.Media)
         addr, err := parse(ref)
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
