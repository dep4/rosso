package json

import (
   "bytes"
   "encoding/json"
)

var (
   Marshal = json.Marshal
   MarshalIndent = json.MarshalIndent
   NewDecoder = json.NewDecoder
   NewEncoder = json.NewEncoder
   Unmarshal = json.Unmarshal
)

type Scanner struct {
   Data []byte
   Sep []byte
}

func (self Scanner) Decode(val any) error {
   data := append(self.Sep, self.Data...)
   dec := NewDecoder(bytes.NewReader(data))
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return json.Unmarshal(data[:high], val)
      }
   }
}

func (self *Scanner) Scan() bool {
   var found bool
   _, self.Data, found = bytes.Cut(self.Data, self.Sep)
   return found
}
