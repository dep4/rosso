package json

import (
   "bytes"
   "encoding/json"
   "github.com/89z/format"
   "os"
)

var (
   NewDecoder = json.NewDecoder
   NewEncoder = json.NewEncoder
)

// `value` is the receiver
func Create[T any](value T, name string) error {
   file, err := format.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(value)
}

func Open[T any](name string) (*T, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   value := new(T)
   if err := json.NewDecoder(file).Decode(value); err != nil {
      return nil, err
   }
   return value, nil
}

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
