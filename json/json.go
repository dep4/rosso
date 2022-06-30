package json

import (
   "bytes"
   "encoding/json"
   "github.com/89z/format"
   "io"
   "os"
)

var (
   Marshal = json.Marshal
   NewDecoder = json.NewDecoder
   NewEncoder = json.NewEncoder
   Unmarshal = json.Unmarshal
)

func Buffer(value any) (*bytes.Buffer, error) {
   buf := new(bytes.Buffer)
   err := indent(buf).Encode(value)
   if err != nil {
      return nil, err
   }
   return buf, nil
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
   return indent(file).Encode(value)
}

func indent(w io.Writer) *json.Encoder {
   enc := json.NewEncoder(w)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc
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
