package dash

import (
   "strconv"
)

func (Representation) String() string {
   return ""
}

func (r Representation) line_3() string {
   b := []byte("MimeType:")
   b = append(b, r.MimeType...)
   if r.Codecs != "" {
      b = append(b, " Codecs:"...)
      b = append(b, r.Codecs...)
   }
   if r.Adaptation.Lang != "" {
      b = append(b, " Lang:"...)
      b = append(b, r.Adaptation.Lang...)
   }
   if r.Adaptation.Role != nil {
      b = append(b, " Role:"...)
      b = append(b, r.Adaptation.Role.Value...)
   }
   return string(b)
}

func (r Representation) line_2() string {
   first := true
   var b []byte
   if r.Width >= 1 {
      b = append(b, "Width:"...)
      b = strconv.AppendInt(b, r.Width, 10)
      b = append(b, " Height:"...)
      b = strconv.AppendInt(b, r.Height, 10)
      first = false
   }
   if r.Bandwidth >= 1 {
      if !first {
         b = append(b, ' ')
      }
      b = append(b, "Bandwidth:"...)
      b = strconv.AppendInt(b, r.Bandwidth, 10)
   }
   return string(b)
}

func (r Representation) line_1() string {
   return "ID:" + r.ID
}

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType string `xml:"mimeType,attr"`
   Width int64 `xml:"width,attr"`
}

type Adaptation struct {
   Codecs string `xml:"codecs,attr"`
   Lang string `xml:"lang,attr"`
   MimeType string `xml:"mimeType,attr"`
   Role *struct {
      Value string `xml:"value,attr"`
   }
}
