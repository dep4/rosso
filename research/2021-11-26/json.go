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
      if buf[0] == '{' {
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
      } else {
         var raw []json.RawMessage
         err := json.Unmarshal(buf, &raw)
         if err != nil {
            return err
         }
         var arr []interface{}
         for _, val := range raw {
            var any token
            err := json.Unmarshal(val, &any)
            if err != nil {
               return err
            }
            arr = append(arr, any)
         }
         (*m)[key] = arr
      }
   }
   return nil
}


type fixed32 uint32

func (f fixed32) MarshalJSON() ([]byte, error) {
   tok := token{
      protowire.Fixed32Type, uint32(f),
   }
   return json.Marshal(tok)
}

func (f *fixed32) UnmarshalJSON(buf []byte) error {
   var tok struct {
      Value uint32
   }
   err := json.Unmarshal(buf, &tok)
   if err != nil {
      return err
   }
   *f = fixed32(tok.Value)
   return nil
}

type fixed64 uint64

func (f fixed64) MarshalJSON() ([]byte, error) {
   tok := token{
      protowire.Fixed64Type, uint64(f),
   }
   return json.Marshal(tok)
}

func (f *fixed64) UnmarshalJSON(buf []byte) error {
   var tok struct {
      Value uint64
   }
   err := json.Unmarshal(buf, &tok)
   if err != nil {
      return err
   }
   *f = fixed64(tok.Value)
   return nil
}

type str string

func (s str) MarshalJSON() ([]byte, error) {
   tok := token{
      protowire.BytesType, string(s),
   }
   return json.Marshal(tok)
}

func (s *str) UnmarshalJSON(buf []byte) error {
   var tok struct {
      Value string
   }
   err := json.Unmarshal(buf, &tok)
   if err != nil {
      return err
   }
   *s = str(tok.Value)
   return nil
}

type varint uint64

func (v varint) MarshalJSON() ([]byte, error) {
   tok := token{
      protowire.VarintType, uint64(v),
   }
   return json.Marshal(tok)
}

func (v *varint) UnmarshalJSON(buf []byte) error {
   var tok struct {
      Value uint64
   }
   err := json.Unmarshal(buf, &tok)
   if err != nil {
      return err
   }
   *v = varint(tok.Value)
   return nil
}
