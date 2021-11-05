package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

var indent int

func consume(n protowire.Number, t protowire.Type, buf []byte) (interface{}, int) {
   switch t {
   case protowire.VarintType:
      return protowire.ConsumeVarint(buf)
   case protowire.Fixed32Type:
      return protowire.ConsumeFixed32(buf)
   case protowire.Fixed64Type:
      return protowire.ConsumeFixed64(buf)
   case protowire.BytesType:
      v, vLen := protowire.ConsumeBytes(buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return string(v), vLen
   case protowire.StartGroupType:
      v, vLen := protowire.ConsumeGroup(n, buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return v, vLen
   }
   return nil, 0
}

type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value interface{}
}

type Fields []Field

func Parse(buf []byte) Fields {
   var fs Fields
   for len(buf) > 0 {
      n, t, fLen := protowire.ConsumeField(buf)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if tLen <= 0 {
         return nil
      }
      v, vLen := consume(n, t, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      fs = append(fs, Field{n, t, v})
      buf = buf[fLen:]
   }
   return fs
}

func (fs Fields) String() string {
   var buf string
   for k, v := range fs {
      if k >= 1 {
         buf += "\n"
      }
      buf += strings.Repeat("   ", indent)
      _, ok := v.Value.(Fields)
      if ok {
         indent++
      }
      buf += fmt.Sprintf("number:%v type:%v value:", v.Number, v.Type)
      if ok {
         buf += "\n"
      }
      buf += fmt.Sprint(v.Value)
   }
   if indent >= 1 {
      indent--
   }
   return buf
}
