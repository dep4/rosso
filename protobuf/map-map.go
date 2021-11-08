package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

type mapJSON = map[string]interface{}

type mapProto map[protowire.Number]interface{}

func newMapProto(m mapJSON) (mapProto, error) {
   pMap := make(mapProto)
   for str, val := range m {
      num, err := strconv.Atoi(str)
      if err != nil {
         return nil, err
      }
      if err := pMap.set(protowire.Number(num), val); err != nil {
         return nil, err
      }
   }
   return pMap, nil
}

func (m mapProto) set(num protowire.Number, any interface{}) error {
   switch val := any.(type) {
   case map[string]interface{}:
      pMap, err := newMapProto(val)
      if err != nil {
         return err
      }
      m[num] = pMap
   case []interface{}:
      for _, any := range val {
         err := m.set(num, any)
         if err != nil {
            return err
         }
      }
   default:
      m[num] = any
   }
   return nil
}
