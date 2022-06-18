package protobuf

import (
   "encoding/base64"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
)

func (t Token) MarshalJSON() ([]byte, error) {
   buf := []byte{'['}
   for key, val := range t.Value {
      if key >= 1 {
         buf = append(buf, ',')
      }
      switch t.Type {
      case protowire.VarintType:
         buf = strconv.AppendUint(buf, val.Int64, 10)
      case protowire.Fixed64Type:
         buf = strconv.AppendUint(buf, val.Int64, 10)
      case protowire.Fixed32Type:
         buf = strconv.AppendUint(buf, uint64(val.Int32), 10)
      case protowire.BytesType:
         var s string
         if format.IsString(val.Bytes) {
            s = string(val.Bytes)
         } else {
            s = base64.StdEncoding.EncodeToString(val.Bytes)
         }
         buf = strconv.AppendQuote(buf, s)
      }
   }
   return append(buf, ']'), nil
}

type Token struct {
   Type protowire.Type
   Value []Value
}

type Value struct {
   Int32 uint32
   Int64 uint64
   Bytes []byte
   Message Message
}

type Message map[protowire.Number]Token

func (m Message) UnmarshalBinary(buf []byte) error {
   if len(buf) == 0 {
      return io.ErrUnexpectedEOF
   }
   for len(buf) >= 1 {
      num, typ, tLen := protowire.ConsumeTag(buf)
      err := protowire.ParseError(tLen)
      if err != nil {
         return err
      }
      buf = buf[tLen:]
      tok := m[num]
      var (
         vLen int
         val Value
      )
      switch typ {
      case protowire.BytesType:
         val.Bytes, vLen = protowire.ConsumeBytes(buf)
         val.Message = make(Message)
         val.Message.UnmarshalBinary(val.Bytes)
      case protowire.VarintType:
         val.Int64, vLen = protowire.ConsumeVarint(buf)
      case protowire.Fixed64Type:
         val.Int64, vLen = protowire.ConsumeFixed64(buf)
      case protowire.Fixed32Type:
         val.Int32, vLen = protowire.ConsumeFixed32(buf)
      }
      tok.Value = append(tok.Value, val)
      m[num] = tok
      if err := protowire.ParseError(vLen); err != nil {
         return err
      }
      buf = buf[vLen:]
   }
   return nil
}
