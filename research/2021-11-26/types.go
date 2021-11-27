package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

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
