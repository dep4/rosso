package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

// google.golang.org/protobuf/encoding/protowire#AppendBytes
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

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (m Message) Get(n protowire.Number, name string) Message {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].(Message)
   if ok {
      return val
   }
   return nil
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (m Message) GetBytes(n protowire.Number, name string) []byte {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].([]byte)
   if ok {
      return val
   }
   return nil
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (m Message) GetMessages(n protowire.Number, name string) []Message {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].([]Message)
   if ok {
      return val
   }
   return nil
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
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

// google.golang.org/protobuf/encoding/protowire#AppendFixed32
func (m Message) addFixed32(n protowire.Number, v uint32) {
   tag := Tag{Type: protowire.Fixed32Type, Number: n}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint32:
      m[tag] = []uint32{val, v}
   case []uint32:
      m[tag] = append(val, v)
   }
}

// google.golang.org/protobuf/encoding/protowire#AppendFixed64
func (m Message) addFixed64(n protowire.Number, v uint64) {
   tag := Tag{Type: protowire.Fixed64Type, Number: n}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint64:
      m[tag] = []uint64{val, v}
   case []uint64:
      m[tag] = append(val, v)
   }
}

// google.golang.org/protobuf/encoding/protowire#AppendGroup
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

// google.golang.org/protobuf/encoding/protowire#AppendString
func (m Message) GetString(n protowire.Number, name string) string {
   tag := Tag{Type: protowire.BytesType, Number: n}
   val, ok := m[tag].(string)
   if ok {
      return val
   }
   return ""
}

// google.golang.org/protobuf/encoding/protowire#AppendString
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

// google.golang.org/protobuf/encoding/protowire#AppendVarint
func (m Message) GetVarint(n protowire.Number, name string) uint64 {
   tag := Tag{Type: protowire.VarintType, Number: n}
   val, ok := m[tag].(uint64)
   if ok {
      return val
   }
   return 0
}

// google.golang.org/protobuf/encoding/protowire#AppendVarint
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

// google.golang.org/protobuf/encoding/protowire#ConsumeBytes
func (m Message) consumeBytes(n protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeBytes(b)
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

// google.golang.org/protobuf/encoding/protowire#ConsumeFixed32
func (m Message) consumeFixed32(n protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addFixed32(n, val)
   return nil
}

// google.golang.org/protobuf/encoding/protowire#ConsumeFixed64
func (m Message) consumeFixed64(n protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addFixed64(n, val)
   return nil
}

// google.golang.org/protobuf/encoding/protowire#ConsumeGroup
func (m Message) consumeGroup(n protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeGroup(n, b)
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

// google.golang.org/protobuf/encoding/protowire#ConsumeVarint
func (m Message) consumeVarint(n protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeVarint(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addVarint(n, val)
   return nil
}
