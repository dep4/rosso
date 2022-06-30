package json

import (
   "bytes"
   "encoding/json"
   "github.com/89z/format"
   "os"
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

func (s Scanner) Decode(val any) error {
   data := append(s.Sep, s.Data...)
   dec := NewDecoder(bytes.NewReader(data))
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return json.Unmarshal(data[:high], val)
      }
   }
}

func (s *Scanner) Scan() bool {
   var found bool
   _, s.Data, found = bytes.Cut(s.Data, s.Sep)
   return found
}

func Decode(name string, value any) error {
   file, err := os.Open(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return json.NewDecoder(file).Decode(value)
}

func Encode(name string, value any) error {
   file, err := format.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc.Encode(value)
}
