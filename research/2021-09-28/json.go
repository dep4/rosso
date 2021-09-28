package js

import (
   "bytes"
   "encoding/json"
)

type scanner struct {
   left []byte
   right []byte
}

func newScanner(b []byte) scanner {
   return scanner{right: b}
}

func (s scanner) bytes() []byte {
   return s.left
}

func (s *scanner) scan() bool {
   for len(s.right) > 0 {
      r := bytes.NewReader(s.right)
      dec := json.NewDecoder(r)
      _, err := dec.Token()
      if err == nil {
         for {
            _, err := dec.Token()
            if err != nil {
               if dec.More() {
                  break
               }
               s.left = s.right[:dec.InputOffset()]
               s.right = s.right[dec.InputOffset():]
               return true
            }
         }
      }
      s.right = s.right[1:]
   }
   return false
}
