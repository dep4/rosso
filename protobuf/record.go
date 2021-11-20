package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

func appendField(buf []byte, num protowire.Number, val interface{}) ([]byte, error) {
   switch val := val.(type) {
   case bool:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, protowire.EncodeBool(val))
   case float64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, uint64(val))
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []interface{}:
      for _, elem := range val {
         aBuf, err := appendField(buf, num, elem)
         if err != nil {
            return nil, err
         }
         buf = aBuf
      }
   case map[string]interface{}:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      eBuf, err := Records.Bytes(val)
      if err != nil {
         return nil, err
      }
      buf = protowire.AppendBytes(buf, eBuf)
   }
   return buf, nil
}

type Records map[string]interface{}

// Convert struct to map
func NewRecords(val interface{}) (Records, error) {
   buf, err := json.Marshal(val)
   if err != nil {
      return nil, err
   }
   var smap Records
   if err := json.Unmarshal(buf, &smap); err != nil {
      return nil, err
   }
   return smap, nil
}

// Convert map to byte slice
func (r Records) Bytes() ([]byte, error) {
   var buf []byte
   for str, val := range r {
      num, err := strconv.Atoi(str)
      if err != nil {
         return nil, err
      }
      buf, err = appendField(buf, protowire.Number(num), val)
      if err != nil {
         return nil, err
      }
   }
   return buf, nil
}
