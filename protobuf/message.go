package protobuf

import (
   "errors"
   "github.com/89z/rosso/strconv"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Message map[Number]Encoder

func Unmarshal(b []byte) (Message, error) {
   if len(b) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(b) >= 1 {
      num, typ, length := protowire.ConsumeTag(b)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      b = b[length:]
      switch typ {
      case protowire.VarintType:
         b, err = mes.consume_varint(num, b)
      case protowire.Fixed64Type:
         b, err = mes.consume_fixed64(num, b)
      case protowire.Fixed32Type:
         b, err = mes.consume_fixed32(num, b)
      case protowire.BytesType:
         b, err = mes.consume_raw(num, b)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}

func (self Message) Add(num Number, value Message) error {
   switch lvalue := self[num].(type) {
   case nil:
      self[num] = value
   case Message:
      self[num] = Slice[Message]{lvalue, value}
   case Slice[Message]:
      self[num] = append(lvalue, value)
   default:
      return type_error{num, lvalue, value}
   }
   return nil
}

func (self Message) Add_Fixed32(num Number, value uint32) error {
   rvalue := Fixed32(value)
   switch lvalue := self[num].(type) {
   case nil:
      self[num] = rvalue
   case Fixed32:
      self[num] = Slice[Fixed32]{lvalue, rvalue}
   case Slice[Fixed32]:
      self[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (self Message) Add_Fixed64(num Number, value uint64) error {
   rvalue := Fixed64(value)
   switch lvalue := self[num].(type) {
   case nil:
      self[num] = rvalue
   case Fixed64:
      self[num] = Slice[Fixed64]{lvalue, rvalue}
   case Slice[Fixed64]:
      self[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (self Message) Add_String(num Number, value string) error {
   rvalue := String(value)
   switch lvalue := self[num].(type) {
   case nil:
      self[num] = rvalue
   case String:
      self[num] = Slice[String]{lvalue, rvalue}
   case Slice[String]:
      self[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (self Message) Add_Varint(num Number, value uint64) error {
   rvalue := Varint(value)
   switch lvalue := self[num].(type) {
   case nil:
      self[num] = rvalue
   case Varint:
      self[num] = Slice[Varint]{lvalue, rvalue}
   case Slice[Varint]:
      self[num] = append(lvalue, rvalue)
   default:
      return type_error{num, lvalue, rvalue}
   }
   return nil
}

func (self Message) Get(num Number) Message {
   switch rvalue := self[num].(type) {
   case Message:
      return rvalue
   case Raw:
      return rvalue.Message
   }
   return nil
}

func (self Message) Get_Bytes(num Number) ([]byte, error) {
   lvalue := self[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return nil, type_error{num, lvalue, rvalue}
   }
   return rvalue.Bytes, nil
}

func (self Message) Get_Fixed64(num Number) (uint64, error) {
   lvalue := self[num]
   rvalue, ok := lvalue.(Fixed64)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (self Message) Get_Messages(num Number) []Message {
   switch rvalue := self[num].(type) {
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

func (self Message) Get_String(num Number) (string, error) {
   lvalue := self[num]
   rvalue, ok := lvalue.(Raw)
   if !ok {
      return "", type_error{num, lvalue, rvalue}
   }
   return rvalue.String, nil
}

func (self Message) Get_Varint(num Number) (uint64, error) {
   lvalue := self[num]
   rvalue, ok := lvalue.(Varint)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (self Message) Marshal() []byte {
   var (
      nums []Number
      bufs []byte
   )
   for num := range self {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      bufs = self[num].encode(bufs, num)
   }
   return bufs
}

func (self Message) consume_fixed32(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := self.Add_Fixed32(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (self Message) consume_fixed64(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := self.Add_Fixed64(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (self Message) consume_raw(num Number, b []byte) ([]byte, error) {
   var (
      length int
      rvalue Raw
   )
   rvalue.Bytes, length = protowire.ConsumeBytes(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if strconv.String(rvalue.Bytes) {
      rvalue.String = string(rvalue.Bytes)
   }
   rvalue.Message, _ = Unmarshal(rvalue.Bytes)
   switch lvalue := self[num].(type) {
   case nil:
      self[num] = rvalue
   case Raw:
      self[num] = Slice[Raw]{lvalue, rvalue}
   case Slice[Raw]:
      self[num] = append(lvalue, rvalue)
   default:
      return nil, type_error{num, lvalue, rvalue}
   }
   return b[length:], nil
}

func (self Message) consume_varint(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeVarint(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := self.Add_Varint(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (self Message) encode(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, self.Marshal())
}

func (Message) get_type() string { return "Message" }
