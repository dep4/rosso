package dash

import (
   "net/url"
   "strconv"
   "strings"
)

func (r Representation) String() string {
   var b []byte
   b = append(b, "Lang:"...)
   b = append(b, r.Adaptation.Lang...)
   if r.Adaptation.Role != nil {
      b = append(b, " Role:"...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   b = append(b, " Bandwidth:"...)
   b = strconv.AppendInt(b, r.Bandwidth, 10)
   b = append(b, " Codecs:"...)
   b = append(b, r.Codecs...)
   if r.Width >= 1 {
      b = append(b, " Width:"...)
      b = strconv.AppendInt(b, r.Width, 10)
      b = append(b, " Height:"...)
      b = strconv.AppendInt(b, r.Height, 10)
   }
   b = append(b, " ID:"...)
   b = append(b, r.ID...)
   return string(b)
}

func (m Media) Representations(fn Represent_Func) Representations {
   var reps Representations
   for _, ada := range m.Period.AdaptationSet {
      for _, rep := range ada.Representation {
         if rep.MIME_Type == "" {
            rep.MIME_Type = ada.MIME_Type
         }
         if rep.SegmentTemplate == nil {
            rep.SegmentTemplate = ada.SegmentTemplate
         }
         if fn(ada, rep) {
            reps = append(reps, rep)
         }
      }
   }
   return reps
}

func (r Representation) Initial(base *url.URL) (*url.URL, error) {
   ref := r.id(r.SegmentTemplate.Initial)
   return base.Parse(ref)
}

func (r Representation) Media(base *url.URL) ([]*url.URL, error) {
   var (
      addrs []*url.URL
      start int
   )
   if r.SegmentTemplate.Start_Number != nil {
      start = *r.SegmentTemplate.Start_Number
   }
   for _, seg := range r.SegmentTemplate.SegmentTimeline.S {
      for seg.T = start; seg.R >= 0; seg.R-- {
         ref := r.id(r.SegmentTemplate.Media)
         if r.SegmentTemplate.Start_Number != nil {
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

func (r Representation) id(in string) string {
   return strings.Replace(in, "$RepresentationID$", r.ID, 1)
}

func (s Segment) number(in string) string {
   return strings.Replace(in, "$Number$", strconv.Itoa(s.T), 1)
}

func (s Segment) time(in string) string {
   return strings.Replace(in, "$Time$", strconv.Itoa(s.T), 1)
}

func (r Representations) Representation(bandwidth int64) *Representation {
   distance := func(r *Representation) int64 {
      if r.Bandwidth > bandwidth {
         return r.Bandwidth - bandwidth
      }
      return bandwidth - r.Bandwidth
   }
   var output *Representation
   for i, input := range r {
      if output == nil || distance(&input) < distance(output) {
         output = &r[i]
      }
   }
   return output
}

