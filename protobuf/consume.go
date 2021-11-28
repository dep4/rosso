package protobuf

import (
   "encoding/base64"
   "google.golang.org/protobuf/encoding/protowire"
)

func consume(num protowire.Number, typ protowire.Type, buf []byte) (interface{}, error) {
   switch typ {
   case protowire.VarintType:
      return consumeVarint(buf)
   case protowire.Fixed32Type:
      return consumeFixed32(buf)
   case protowire.Fixed64Type:
      return consumeFixed64(buf)
   case protowire.BytesType:
      val, err := consumeBytes(buf)
      if err != nil {
         return nil, err
      }
      mes, err := Unmarshal(val)
      if err != nil {
         if isBinary(val) {
            return base64.StdEncoding.EncodeToString(val), nil
         }
         return string(val), nil
      }
      return mes, nil
   case protowire.StartGroupType:
      val, err := consumeGroup(num, buf)
      if err != nil {
         return nil, err
      }
      return Unmarshal(val)
   }
   return nil, nil
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
