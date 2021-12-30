package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

type Message map[Tag]interface{}

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen, err := consumeField(buf)
      if err != nil {
         return nil, err
      }
      tLen, err := consumeTag(buf[:fLen])
      if err != nil {
         return nil, err
      }
      bVal := buf[tLen:fLen]
      switch typ {
      case protowire.VarintType:
         val, err := consumeVarint(bVal)
         if err != nil {
            return nil, err
         }
         mes.addUint64(num, val)
      case protowire.Fixed64Type:
         val, err := consumeFixed64(bVal)
         if err != nil {
            return nil, err
         }
         mes.addUint64(num, val)
      case protowire.Fixed32Type:
         val, err := consumeFixed32(bVal)
         if err != nil {
            return nil, err
         }
         mes.addUint32(num, val)
      case protowire.BytesType:
         buf, err := consumeBytes(bVal)
         if err != nil {
            return nil, err
         }
         if !isBinary(buf) {
            mes.addString(num, string(buf))
         } else {
            mNew, err := Unmarshal(buf)
            if err != nil {
               mes.addBytes(num, buf)
            } else {
               mes.add(num, mNew)
            }
         }
         /*
         mNew, err := Unmarshal(buf)
         if err != nil {
            if isBinary(buf) {
               mes.addBytes(num, buf)
            } else {
               mes.addString(num, string(buf))
            }
         } else {
            mes.add(num, mNew)
         }
         */
      case protowire.StartGroupType:
         buf, err := consumeGroup(num, bVal)
         if err != nil {
            return nil, err
         }
         mNew, err := Unmarshal(buf)
         if err != nil {
            return nil, err
         }
         mes.add(num, mNew)
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func isBinary(buf []byte) bool {
   for _, b := range buf {
      switch {
      case b <= 0x08,
      b == 0x0B,
      0x0E <= b && b <= 0x1A,
      0x1C <= b && b <= 0x1F:
         return true
      }
   }
   return false
}

func (m Message) add(key protowire.Number, val Message) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{typ, val}
   case []Message:
      m[tag] = append(typ, val)
   }
}

type Tag struct {
   protowire.Number
   String string
}

func consumeBytes(b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeBytes(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return nil, err
   }
   return val, nil
}

func consumeField(b []byte) (protowire.Number, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(b)
   err := protowire.ParseError(fLen)
   if err != nil {
      return 0, 0, 0, err
   }
   return num, typ, fLen, nil
}

func consumeFixed32(b []byte) (uint32, error) {
   val, vLen := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeFixed64(b []byte) (uint64, error) {
   val, vLen := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeGroup(num protowire.Number, b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeGroup(num, b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return nil, err
   }
   return val, nil
}

func consumeTag(b []byte) (int, error) {
   _, _, tLen := protowire.ConsumeTag(b)
   err := protowire.ParseError(tLen)
   if err != nil {
      return 0, err
   }
   return tLen, nil
}

func consumeVarint(b []byte) (uint64, error) {
   val, vLen := protowire.ConsumeVarint(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func (m Message) addBytes(key protowire.Number, val []byte) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case []byte:
      m[tag] = [][]byte{typ, val}
   case [][]byte:
      m[tag] = append(typ, val)
   }
}

func (m Message) addString(key protowire.Number, val string) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{typ, val}
   case []string:
      m[tag] = append(typ, val)
   }
}

func (m Message) addUint32(key protowire.Number, val uint32) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint32:
      m[tag] = []uint32{typ, val}
   case []uint32:
      m[tag] = append(typ, val)
   }
}

func (m Message) addUint64(key protowire.Number, val uint64) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{typ, val}
   case []uint64:
      m[tag] = append(typ, val)
   }
}
