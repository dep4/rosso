package json

import (
   "bytes"
   "io"
)

type Scanner struct {
   Split []byte
   buf []byte
}

func NewScanner(r io.Reader) (*Scanner, error) {
   buf, err := io.ReadAll(r)
   if err != nil {
      return nil, err
   }
   return &Scanner{buf: buf}, nil
}

func (s *Scanner) ReadFrom(r io.Reader) (int64, error) {
   var buf bytes.Buffer
   num, err := io.Copy(&buf, r)
   if err != nil {
      return 0, err
   }
   s.buf = buf.Bytes()
   return num, nil
}
