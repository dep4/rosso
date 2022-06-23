package protobuf

import (
   "errors"
   "github.com/89z/format"
   "io"
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Number]Encoder

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, length := protowire.ConsumeTag(buf)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      buf = buf[length:]
      switch typ {
      case protowire.VarintType:
         buf, err = mes.consume_varint(num, buf)
      case protowire.Fixed64Type:
         buf, err = mes.consume_fixed64(num, buf)
      case protowire.Fixed32Type:
         buf, err = mes.consume_fixed32(num, buf)
      case protowire.BytesType:
         buf, err = mes.consume_raw(num, buf)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}

func (m Message) Add_Fixed32(num Number, val uint32) error {
   rvalue := Fixed32(val)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Fixed32:
      m[num] = Slice[Fixed32]{lvalue, rvalue}
   case Slice[Fixed32]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Add_Fixed64(num Number, val uint64) error {
   rvalue := Fixed64(val)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Fixed64:
      m[num] = Slice[Fixed64]{lvalue, rvalue}
   case Slice[Fixed64]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Add_Varint(num Number, val uint64) error {
   rvalue := Varint(val)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Varint:
      m[num] = Slice[Varint]{lvalue, rvalue}
   case Slice[Varint]:
      m[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) Get(num Number) Message {
   switch rvalue := m[num].(type) {
   case Message:
      return rvalue
   case Raw:
      return rvalue.Message
   }
   return nil
}

func (m Message) Get_Bytes(num Number) ([]byte, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return nil, type_error{num, lvalue, rvalue}
   }
   return rvalue.Bytes, nil
}

func (m Message) Get_Fixed64(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Fixed64)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) Get_Messages(num Number) []Message {
   switch rvalue := m[num].(type) {
   case Message:
      return []Message{rvalue}
   case Slice[Message]:
      return rvalue
   case Raw:
      return []Message{rvalue.Message}
   case Slice[Raw]:
      var mes []Message
      for _, raw := range rvalue {
         mes = append(mes, raw.Message)
      }
      return mes
   }
   return nil
}

func (m Message) Get_String(num Number) (string, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return "", type_error{num, lvalue, rvalue}
   }
   return rvalue.String, nil
}

func (m Message) Get_Varint(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Varint)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) consume_fixed32(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Fixed32(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (m Message) consume_fixed64(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Fixed64(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (m Message) consume_raw(num Number, b []byte) ([]byte, error) {
   var (
      length int
      rvalue Raw
   )
   rvalue.Bytes, length = protowire.ConsumeBytes(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if format.String(rvalue.Bytes) {
      rvalue.String = string(rvalue.Bytes)
   }
   rvalue.Message, _ = Unmarshal(rvalue.Bytes)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Raw:
      m[num] = Slice[Raw]{lvalue, rvalue}
   case Slice[Raw]:
      m[num] = append(lvalue, rvalue)
   default:
      return nil, type_error{num, lvalue, rvalue}
   }
   return b[length:], nil
}

func (m Message) consume_varint(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeVarint(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Varint(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (Message) get_type() string { return "Message" }
