package protobuf

import (
   "encoding/json"
)

func toMap(any interface{}) (map[string]interface{}, error) {
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

func (m Message) toStruct(any interface{}) error {
   buf, err := json.Marshal(m)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, any)
}
