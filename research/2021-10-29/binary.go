package binary

import (
   "encoding/binary"
)

type decoder struct {
   data []byte
   low, high int
}

func newDecoder(data []byte) *decoder {
   return &decoder{
      data: data, high: len(data),
   }
}

func (d decoder) slice() []byte {
   if d.high < d.low {
      return nil
   }
   if d.high < 0 {
      return nil
   }
   if d.high > len(d.data) {
      return nil
   }
   return d.data[d.low:d.high]
}

func (d *decoder) uint16LengthPrefixed() []byte {
   d.high = d.low + 2
   data := d.slice()
   if data == nil {
      return nil
   }
   d.low += 2
   d.high = d.low + int(binary.BigEndian.Uint16(data))
   data = d.slice()
   d.low = d.high
   return data
}

func (d *decoder) uint32LengthPrefixed() []byte {
   d.high = d.low + 4
   data := d.slice()
   if data == nil {
      return nil
   }
   d.low += 4
   d.high = d.low + int(binary.BigEndian.Uint32(data))
   data = d.slice()
   d.low = d.high
   return data
}

func (d *decoder) seek(n int) {
   d.low += n
}

func (d *decoder) seekByte(b byte) {
   d.low += bytes.IndexByte(d.slice(), b)
}

func handshakes(data []byte) [][]byte {
   var hands [][]byte
   dec := newDecoder(data)
   for {
      dec.seekByte(0x16)
      // Content Type
      dec.seek(1)
      // Version
      dec.seek(2)
      // Length, Handshake Protocol
      dec.uint16LengthPrefixed()
      // FIXME
      dec.low = low
      hand := dec.slice()
      if len(data) > 0 {
         hands = append(hands, hand)
      }
      dec.seek(1)
   }
}
