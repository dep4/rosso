package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

// mimesniff.spec.whatwg.org#binary-data-byte
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

type Tag struct {
   protowire.Number
   string
}

type Message map[Tag]interface{}

func Decode(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

func (m Message) Encode() io.Reader {
   buf := m.Marshal()
   return bytes.NewReader(buf)
}

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

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case uint32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
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
      for _, ran := range val {
         buf = appendField(buf, num, ran)
      }
   case []string:
      for _, ran := range val {
         buf = appendField(buf, num, ran)
      }
   case []Message:
      for _, ran := range val {
         buf = appendField(buf, num, ran)
      }
   }
   return buf
}

func (m Message) Marshal() []byte {
   var buf []byte
   for key, val := range m {
      buf = appendField(buf, key.Number, val)
   }
   return buf
}
