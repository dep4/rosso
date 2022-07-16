package xml

import (
   "bytes"
   "encoding/xml"
)

func decoder(data []byte) *xml.Decoder {
   dec := xml.NewDecoder(bytes.NewReader(data))
   dec.AutoClose = xml.HTMLAutoClose
   dec.Strict = false
   return dec
}

type Scanner struct {
   Data []byte
   Sep []byte
}

func (self Scanner) Decode(val any) error {
   data := append(self.Sep, self.Data...)
   dec := decoder(data)
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return decoder(data[:high]).Decode(val)
      }
   }
}

func (self *Scanner) Scan() bool {
   var found bool
   _, self.Data, found = bytes.Cut(self.Data, self.Sep)
   return found
}
