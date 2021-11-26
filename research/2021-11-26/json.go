package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m message) MarshalJSON() ([]byte, error) {
   mes := map[protowire.Number]interface{}(m)
   return json.Marshal(mes)
}

func (m *message) UnmarshalJSON(buf []byte) error {
   var raw map[protowire.Number]json.RawMessage
   err := json.Unmarshal(buf, &raw)
   if err != nil {
      return err
   }
   for _, buf := range raw {
      var raw struct {
         Type protowire.Type
         Value json.RawMessage
      }
      err := json.Unmarshal(buf, &raw)
      if err != nil {
         return err
      }
      switch raw.Type {
      case protowire.VarintType:
         var val uint64
         err := json.Unmarshal(raw.Value, &val)
         if err != nil {
            return err
         }
      case protowire.Fixed32Type:
         var val uint32
         err := json.Unmarshal(raw.Value, &val)
         if err != nil {
            return err
         }
      case protowire.Fixed64Type:
         var val uint64
         err := json.Unmarshal(raw.Value, &val)
         if err != nil {
            return err
         }
      case protowire.BytesType:
      }
   }
   return nil
}
