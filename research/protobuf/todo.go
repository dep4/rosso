package protobuf

import (
   "bufio"
   "bytes"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func add[T Encoder](mes Message, num Number, val T) bool {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Encoders[T]{value, val}
   case Encoders[T]:
      mes[num] = append(value, val)
   default:
      return false
   }
   return true
}

func Decode(r io.Reader) (Message, error) {
   buf := bufio.NewReader(r)
   mes := make(Message)
   for {
      num, typ, err := consumeTag(buf)
      if err == io.EOF {
         break
      } else if err != nil {
         return nil, err
      }
      switch typ {
      case protowire.EndGroupType:
         break
      case protowire.VarintType: // 0
         val, err := binary.ReadUvarint(buf)
         if err != nil {
            return nil, err
         }
         add(mes, num, Varint(val))
      case protowire.Fixed64Type: // 1
         var val Fixed64
         err := binary.Read(buf, binary.LittleEndian, &val)
         if err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.Fixed32Type: // 5
         var val Fixed32
         err := binary.Read(buf, binary.LittleEndian, &val)
         if err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.BytesType:
         var val Bytes
         val.Raw, err = consumeBytes(buf)
         if err != nil {
            return nil, err
         }
         val.Message, _ = Decode(bytes.NewReader(val.Raw))
         add(mes, num, val)
      case protowire.StartGroupType:
         val, err := Decode(buf)
         if err != nil {
            return nil, err
         }
         add(mes, num, val)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
   }
   return mes, nil
}
