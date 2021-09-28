package json

import (
   "bytes"
   "encoding/json"
)

type Scanner struct {
   left []byte
   right []byte
}

func NewScanner(b []byte) Scanner {
   return Scanner{right: b}
}

func (s Scanner) Bytes() []byte {
   return s.left
}

func (s *Scanner) Scan() bool {
   for len(s.right) > 0 {
      r := bytes.NewReader(s.right)
      dec := json.NewDecoder(r)
      _, err := dec.Token()
      if err == nil {
         for {
            _, err := dec.Token()
            if err != nil {
               s.left = s.right[:dec.InputOffset()]
               s.right = s.right[dec.InputOffset():]
               if json.Valid(s.left) {
                  return true
               }
               break
            }
         }
      }
      s.right = s.right[1:]
   }
   return false
}
