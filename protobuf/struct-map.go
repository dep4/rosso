package protobuf

import (
   "bytes"
   "encoding/json"
)

func transcode(v interface{}) (map[string]interface{}, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(v)
   if err != nil {
      return nil, err
   }
   var m map[string]interface{}
   if err := json.NewDecoder(buf).Decode(&m); err != nil {
      return nil, err
   }
   return m, nil
}
