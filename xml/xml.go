package xml

import (
   "bytes"
   "encoding/xml"
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

func (s Scanner) Decode(val any) error {
   buf := append(s.Split, s.buf...)
   dec := xml.NewDecoder(bytes.NewReader(buf))
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return xml.Unmarshal(buf[:high], val)
      }
   }
}

func (s *Scanner) Scan() bool {
   var found bool
   _, s.buf, found = bytes.Cut(s.buf, s.Split)
   return found
}