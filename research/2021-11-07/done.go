package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

func appendField(buf []byte, num protowire.Number, val interface{}) ([]byte, error) {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   switch val := val.(type) {
   case map[string]interface{}:
      mar, err := marshal(val)
      if err != nil {
         return nil, err
      }
      buf = protowire.AppendBytes(buf, mar)
   case string:
      buf = protowire.AppendString(buf, val)
   }
   return buf, nil
}

func marshal(m map[string]interface{}) ([]byte, error) {
   var buf []byte
   for str, val := range m {
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
