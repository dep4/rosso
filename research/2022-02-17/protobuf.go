package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

const (
   bytesType = "bytes"
   fixed64Type = "fixed64"
   messageType = "message"
   stringType = "string"
   varintType = "varint"
)

func (m Message) Add(num protowire.Number, name string, val Message) error {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{value, val}
   case []Message:
      m[tag] = append(value, val)
   }
   return nil
}

func (m Message) addString(num protowire.Number, val string) {
   tag := Tag{num, stringType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{value, val}
   case []string:
      m[tag] = append(value, val)
   }
}

func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(elem2)
   mes, err := Unmarshal(elem2)
   if err != nil {
      if binary {
         tag := Tag{num, bytesType}
         switch elem := m[tag].(type) {
         case nil:
            m[tag] = elem2
         case []byte:
            m[tag] = [][]byte{elem, elem2}
         case [][]byte:
            m[tag] = append(elem, elem2)
         }
      } else {
         m.addString(num, string(elem2))
      }
   } else {
      // In this section, if input is binary, then result could be a Message or
      // []byte. We assume for now its always a Message. If input is not binary,
      // then result could be a Message or string. Since its not possible to
      // tell Message from string, we just add both under the same Number, each
      // with its own Type.
      m.Add(num, "", mes)
      if !binary {
         m.addString(num, string(elem2))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   tag := Tag{num, fixed64Type}
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

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   tag := Tag{num, varintType}
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

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []byte:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val)
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case []uint64:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []string:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []Message:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   }
   return buf
}

type Message map[Tag]interface{}

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      err := protowire.ParseError(fLen)
      if err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      val := buf[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         err = mes.consumeBytes(num, val)
      case protowire.Fixed64Type:
         err = mes.consumeFixed64(num, val)
      case protowire.VarintType:
         err = mes.consumeVarint(num, val)
      }
      if err != nil {
         return nil, err
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Marshal() []byte {
   var buf []byte
   for tag, val := range m {
      buf = appendField(buf, tag.Number, val)
   }
   return buf
}

type Tag struct {
   protowire.Number
   Name string
}


