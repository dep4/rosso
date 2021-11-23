package json

import (
   "bytes"
   "encoding/json"
)

type Decoder struct {
   buf []byte
}

func NewDecoder(buf []byte) *Decoder {
   return &Decoder{buf}
}

func (d *Decoder) Decode(val interface{}, c byte) bool {
   for {
      off := bytes.IndexByte(d.buf, c)
      if off == -1 {
         return false
      }
      d.buf = d.buf[off:]
      dec := json.NewDecoder(bytes.NewReader(d.buf))
      for {
         _, err := dec.Token()
         if err != nil {
            off := dec.InputOffset()
            err := json.Unmarshal(d.buf[:off], val)
            d.buf = d.buf[1:]
            if err == nil {
               return true
            }
            break
         }
      }
   }
}

func (d *Decoder) Object(val interface{}) bool {
   return d.Decode(val, '{')
}
