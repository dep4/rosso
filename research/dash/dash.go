package dash

import (
   "encoding/xml"
   "io"
)

type Media struct {
   Period struct {
      AdaptationSet []struct {
         Lang string `xml:"lang,attr"`
         MimeType string `xml:"mimeType,attr"`
         Role *struct {
            Value string `xml:"value,attr"`
         }
      }
   }
}

func (m *Media) ReadFrom(r io.Reader) (int64, error) {
   data, err := io.ReadAll(r)
   if err != nil {
      return 0, err
   }
   if err := xml.Unmarshal(data, m); err != nil {
      return 0, err
   }
   return len(data), nil
}
