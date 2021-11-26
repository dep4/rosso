package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

type varint uint64

func (v varint) MarshalJSON() ([]byte, error) {
   tok := token{
      protowire.VarintType, uint64(v),
   }
   return json.Marshal(tok)
}
