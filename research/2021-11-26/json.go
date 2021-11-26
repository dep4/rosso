package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

func consumeJSON(buf []byte) (token, error) {
   return token{}, nil
}

func (m *message) UnmarshalJSON(buf []byte) error {
   var raw map[protowire.Number]json.RawMessage
   err := json.Unmarshal(buf, &raw)
   if err != nil {
      return err
   }
   for num, val := range raw {
      any, err := consumeJSON(val)
      if err != nil {
         return err
      }
      (*m)[num] = any
   }
   return nil
}

func (m message) MarshalJSON() ([]byte, error) {
   mes := map[protowire.Number]interface{}(m)
   return json.Marshal(mes)
}
