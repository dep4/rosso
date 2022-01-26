package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) consumeBytes(n protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   ok := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if ok {
         m.addBytes(n, val)
      } else {
         m.addString(n, string(val))
      }
   } else if ok {
      // Could be Message or []byte
      m.Add(n, "", mes)
   } else {
      // Cound be Message or string
      m.addString(n, string(val))
   }
   return nil
}

func (m Message) consumeGroup(n protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeGroup(n, buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   mes, err := Unmarshal(val)
   if err != nil {
      return err
   }
   m.addGroup(n, mes)
   return nil
}

func (m Message) consumeVarint(n protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addVarint(n, val)
   return nil
}

func (m Message) GetVarint(n protowire.Number, name string) uint64 {
   tag := Tag{Type: protowire.VarintType, Number: n}
   val, ok := m[tag].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) Add(n protowire.Number, name string, v Message) {
   tag := Tag{protowire.BytesType, n, name}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case Message:
      m[tag] = []Message{val, v}
   case []Message:
      m[tag] = append(val, v)
   }
}

func (m Message) Get(n protowire.Number, name string) Message {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].(Message)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetBytes(n protowire.Number, name string) []byte {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].([]byte)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetString(n protowire.Number, name string) string {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].(string)
   if ok {
      return val
   }
   return ""
}

func (m Message) GetMessages(n protowire.Number, name string) []Message {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].([]Message)
   if ok {
      return val
   }
   return nil
}

func (m Message) consumeFixed32(n protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   tag := Tag{Type: protowire.Fixed32Type, Number: n}
   switch elem := m[tag].(type) {
   case nil:
      m[tag] = elem2
   case uint32:
      m[tag] = []uint32{elem, elem2}
   case []uint32:
      m[tag] = append(elem, elem2)
   }
   return nil
}

func (m Message) consumeFixed64(n protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   tag := Tag{Type: protowire.Fixed64Type, Number: n}
   switch elem := m[tag].(type) {
   case nil:
      m[tag] = elem2
   case uint64:
      m[tag] = []uint64{elem, elem2}
   case []uint64:
      m[tag] = append(elem, elem2)
   }
   return nil
}

////////////////////////////////////////////////////////////////////////////////

func (m Message) addVarint(n protowire.Number, v uint64) {
   tag := Tag{Type: protowire.VarintType, Number: n}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint64:
      m[tag] = []uint64{val, v}
   case []uint64:
      m[tag] = append(val, v)
   }
}

func (m Message) addBytes(n protowire.Number, v []byte) {
   tag := Tag{Type: protowire.BytesType, Number: n}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case []byte:
      m[tag] = [][]byte{val, v}
   case [][]byte:
      m[tag] = append(val, v)
   }
}

func (m Message) addGroup(n protowire.Number, v Message) {
   tag := Tag{Type: protowire.StartGroupType, Number: n}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case Message:
      m[tag] = []Message{val, v}
   case []Message:
      m[tag] = append(val, v)
   }
}

func (m Message) addString(n protowire.Number, v string) {
   tag := Tag{Type: protowire.BytesType, Number: n}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case string:
      m[tag] = []string{val, v}
   case []string:
      m[tag] = append(val, v)
   }
}
