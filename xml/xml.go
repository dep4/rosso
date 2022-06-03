package xml

import (
   "bytes"
   "encoding/xml"
   "io"
)

func NewDecoder(buf []byte) *xml.Decoder {
   src := bytes.NewReader(buf)
   dec := xml.NewDecoder(src)
   dec.AutoClose = xml.HTMLAutoClose
   dec.Strict = false
   return dec
}

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
   dec := NewDecoder(buf)
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return NewDecoder(buf[:high]).Decode(val)
      }
   }
}

func (s *Scanner) Scan() bool {
   var found bool
   _, s.buf, found = bytes.Cut(s.buf, s.Split)
   return found
}
