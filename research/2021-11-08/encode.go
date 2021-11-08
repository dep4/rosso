package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

func appendOld(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case Repeated:
      for _, v := range val {
         buf = appendField(buf, num, v)
      }
   case bool:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, protowire.EncodeBool(val))
   case int32:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, uint64(val))
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   }
   return buf
}

func appendNew(buf []byte, num protowire.Number, val interface{}) ([]byte, error) {
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

func alfa(any interface{}) (map[string]interface{}, error) {
   buf, err := json.Marshal(any)
   if err != nil {
      return nil, err
   }
   var mJSON map[string]interface{}
   if err := json.Unmarshal(buf, &mJSON); err != nil {
      return nil, err
   }
   return mJSON, nil
}

func bravo(m map[string]interface{}) ([]byte, error) {
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
