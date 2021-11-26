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
   for key, buf := range raw {
      var raw struct {
         Type protowire.Type
         Value json.RawMessage
      }
      err := json.Unmarshal(buf, &raw)
      if err != nil {
         return err
      }
      switch raw.Type {
      case protowire.Fixed32Type:
         var val uint32
         err := json.Unmarshal(raw.Value, &val)
         if err != nil {
            return err
         }
         (*m)[key] = token{raw.Type, val}
      case protowire.Fixed64Type, protowire.VarintType:
         var val uint64
         err := json.Unmarshal(raw.Value, &val)
         if err != nil {
            return err
         }
         (*m)[key] = token{raw.Type, val}
      case protowire.BytesType:
         if raw.Value[0] == '"' {
            var val string
            err := json.Unmarshal(raw.Value, &val)
            if err != nil {
               return err
            }
            (*m)[key] = token{raw.Type, val}
         } else {
            val := make(message)
            err := json.Unmarshal(raw.Value, &val)
            if err != nil {
               return err
            }
            (*m)[key] = token{raw.Type, val}
         }
      }
   }
   return nil
}
