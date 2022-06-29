package dash

type Representations []Representation

type Representation struct {
   Adaptation *Adaptation
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   ContentProtection *ContentProtection
   SegmentTemplate *SegmentTemplate
   Height int64 `xml:"height,attr"`
   Width int64 `xml:"width,attr"`
   MIME_Type string `xml:"mimeType,attr"`
   ID string `xml:"id,attr"`
}

type Adaptation struct {
   ContentProtection *ContentProtection
   Lang string `xml:"lang,attr"` // Adaptation only
   MIME_Type string `xml:"mimeType,attr"`
   Representation Representations
   SegmentTemplate *SegmentTemplate
   Role *struct { // Adaptation only
      Value string `xml:"value,attr"`
   }
}

type Media struct {
   Period struct {
      AdaptationSet []Adaptation
   }
}

type ContentProtection struct {
   Default_KID string `xml:"default_KID,attr"`
}

type SegmentTemplate struct {
   Initial string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []struct {
         D int `xml:"d,attr"`
         R int `xml:"r,attr"`
         T int `xml:"t,attr"`
      }
   }
   Start_Number *int `xml:"startNumber,attr"`
}
