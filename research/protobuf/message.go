package protobuf

import (
   "errors"
   "github.com/89z/format"
   "io"
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

type Raw struct {
   Bytes []byte
   String string
   Message map[Number]bravo
}

func (Raw) alfa() string { return "Raw" }

type Bytes []byte

func (Bytes) alfa() string { return "Bytes" }

type Fixed32 uint32

func (Fixed32) alfa() string { return "Fixed32" }

type Fixed64 uint64

func (Fixed64) alfa() string { return "Fixed64" }

type Number = protowire.Number

type alfa interface {
   alfa() string
}

type bravo interface {
   bravo() string
}

type String string

func (String) alfa() string { return "String" }

type Varint uint64

func (Varint) alfa() string { return "Varint" }

type type_error struct {
   Number
   lvalue bravo
   rvalue bravo
}

type Message map[Number]bravo

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

func (m Message) Bytes(num Number) ([]byte, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return nil, type_error{num, lvalue, rvalue}
   }
   return rvalue.Bytes, nil
}

func (m Message) Fixed64(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Fixed64)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) String(num Number) (string, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return "", type_error{num, lvalue, rvalue}
   }
   return rvalue.String, nil
}

func (m Message) Varint(num Number) (uint64, error) {
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

func (Message) alfa() string { return "Message" }

type Slice[T alfa] []T

func (Slice[T]) bravo() string {
   var value T
   return "[]" + value.alfa()
}

func (t type_error) Error() string {
   var buf strings.Builder
   buf.WriteString("lvalue ")
   buf.WriteString(t.lvalue.bravo())
   buf.WriteString(" rvalue ")
   buf.WriteString(t.rvalue.bravo())
   return buf.String()
}

////////////////////////////////////////////////////

func (m Message) Message(num Number) Message {
   switch rvalue := m[num].(type) {
   case Message:
      return rvalue
   case Raw:
      return rvalue.Message
   }
   return nil
}

func (m Message) Messages(num Number) []Message {
   var mes []Message
   switch rvalue := m[num].(type) {
   case Raw:
      mes = append(mes, rvalue.Message)
   case Slice[Raw]:
      for _, raw := range rvalue {
         mes = append(mes, raw.Message)
      }
   }
   return mes
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
