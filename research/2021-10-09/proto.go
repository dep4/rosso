package proto

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type token struct {
   Tag struct {
      Num protowire.Number
      Typ protowire.Type
   }
   Length int
   Tokens interface{}
}

func parseUnknown(b []byte) []token {
   var toks []token
   for len(b) > 0 {
      n, t, fieldlen := protowire.ConsumeField(b)
      if fieldlen < 1 {
         return nil
      }
      var tok token
      tok.Tag.Num = n
      tok.Tag.Typ = t
      _, _, taglen := protowire.ConsumeTag(b[:fieldlen])
      if taglen < 1 {
         return nil
      }
      var (
         v interface{}
         vlen int
      )
      switch t {
      case protowire.VarintType:
         v, vlen = protowire.ConsumeVarint(b[taglen:fieldlen])
      case protowire.Fixed32Type:
         v, vlen = protowire.ConsumeFixed32(b[taglen:fieldlen])
      case protowire.Fixed64Type:
         v, vlen = protowire.ConsumeFixed64(b[taglen:fieldlen])
      case protowire.BytesType:
         v, vlen = protowire.ConsumeBytes(b[taglen:fieldlen])
         sub := parseUnknown(v.([]byte))
         if sub != nil {
            v = sub
         }
      case protowire.StartGroupType:
         v, vlen = protowire.ConsumeGroup(n, b[taglen:fieldlen])
         sub := parseUnknown(v.([]byte))
         if sub != nil {
            v = sub
         }
      }
      if vlen < 1 {
         return nil
      }
      tok.Length = vlen - taglen
      bytes, ok := v.([]byte)
      if ok {
         tok.Tokens = string(bytes)
      } else {
         tok.Tokens = v
      }
      toks = append(toks, tok)
      b = b[fieldlen:]
   }
   return toks
}
