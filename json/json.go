package json

import (
   "bytes"
   "encoding/json"
   "io"
)

type Scanner struct {
   Split []byte
   buf []byte
}

func NewScanner(src io.Reader) (*Scanner, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return &Scanner{buf: buf}, nil
}

func (s Scanner) Bytes() []byte {
   return append(s.Split, s.buf...)
}

func (s Scanner) Decode(val any) error {
   buf := s.Bytes()
   dec := json.NewDecoder(bytes.NewReader(buf))
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return json.Unmarshal(buf[:high], val)
      }
   }
}

func (s *Scanner) Scan() bool {
   var found bool
   _, s.buf, found = bytes.Cut(s.buf, s.Split)
   return found
}
