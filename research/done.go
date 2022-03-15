package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

const messageType = 6

type token struct {
   protowire.Number
   protowire.Type
   value interface{}
}

type message []token

func (m message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   for i, tok := range m {
      if tok.Number == num {
         switch value := tok.value.(type) {
         case uint64:
            m[i].value = []uint64{value, val}
         case []uint64:
            m[i].value = append(value, val)
         }
         return nil
      }
   }
   m = append(m, token{num, protowire.Fixed64Type, val})
   return nil
}
