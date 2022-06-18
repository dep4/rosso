package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) consumeVarint(num Number, b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeVarint(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   if in := m[num]; in == nil {
      m[num] = SliceVarint{val}
   } else if out, ok := in.(SliceVarint); ok {
      m[num] = append(out, val)
   } else {
      return nil, typeError{num, in, out}
   }
   return b[vLen:], nil
}

func (m Message) consumeBytes(num Number, b []byte) ([]byte, error) {
   /*
   var val Bytes
   val.Message = make(Message)
   val.Raw, vLen = protowire.ConsumeBytes(buf)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   err := val.Message.UnmarshalBinary(val.Raw)
   if err != nil {
      val.Message = nil
   }
   mes[num] = append(mes[num], val)
   */
   val, vLen := protowire.ConsumeBytes(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   if in := m[num]; in == nil {
      m[num] = SliceBytes{Bytes{Raw: val}}
   } else if out, ok := in.(SliceBytes); ok {
      m[num] = append(out, Bytes{Raw: val})
   } else {
      return nil, typeError{num, in, out}
   }
   return b[vLen:], nil
}

func (m Message) consumeFixed32(num Number, b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeFixed32(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   if in := m[num]; in == nil {
      m[num] = SliceFixed32{val}
   } else if out, ok := in.(SliceFixed32); ok {
      m[num] = append(out, val)
   } else {
      return nil, typeError{num, in, out}
   }
   return b[vLen:], nil
}

func (m Message) consumeFixed64(num Number, b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeFixed64(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   if in := m[num]; in == nil {
      m[num] = SliceFixed64{val}
   } else if out, ok := in.(SliceFixed64); ok {
      m[num] = append(out, val)
   } else {
      return nil, typeError{num, in, out}
   }
   return b[vLen:], nil
}
