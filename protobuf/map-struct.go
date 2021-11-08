package protobuf

import (
   "bytes"
   "encoding/json"
)

func (m Message) Transcode(v interface{}) error {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(m)
   if err != nil {
      return err
   }
   return json.NewDecoder(buf).Decode(v)
}
