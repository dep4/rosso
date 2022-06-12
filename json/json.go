package json

import (
   "bytes"
   "encoding/json"
   "io"
)

// this uses less allocations than `io.ReadAll`
func (s *Scanner) ReadFrom(r io.Reader) (int64, error) {
   var buf bytes.Buffer
   num, err := buf.ReadFrom(r)
   if err != nil {
      return 0, err
   }
   s.buf = buf.Bytes()
   return num, nil
}

var (
   NewDecoder = json.NewDecoder
   NewEncoder = json.NewEncoder
)

type Scanner struct {
   Split []byte
   buf []byte
}

func (s Scanner) Decode(val any) error {
   buf := append(s.Split, s.buf...)
   dec := NewDecoder(bytes.NewReader(buf))
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
