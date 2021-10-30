package binary

import (
   "bytes"
   "encoding/binary"
)

// L1        L2 H2        H3                H1
//  |         |  |         |                 |
//  +---------+--+---------+-----------------+
func handshakes(data []byte) [][]byte {
   var hands [][]byte
   for {
      L1 := bytes.IndexByte(data, 0x16)
      if L1 == -1 {
         return hands
      }
      L2 := L1
      // skip content type
      L2 += 1
      // skip version
      L2 += 2
      H2 := L2
      // skip length
      H2 += 2
      H3 := H2
      // skip handshake
      length := slice(data, L2, H2)
      if length != nil {
         H3 += int(binary.BigEndian.Uint16(length))
         hand := slice(data, L1, H3)
         if hand != nil {
            hands = append(hands, hand)
         }
      }
      data = data[1:]
   }
}

func slice(data []byte, low, high int) []byte {
   if high <= low {
      return nil
   }
   if high < 0 {
      return nil
   }
   if high > len(data) {
      return nil
   }
   return data[low:high]
}
