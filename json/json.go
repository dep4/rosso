package json

import (
   "bytes"
   "encoding/json"
   "fmt"
)

func UnmarshalArray(data []byte, v interface{}) error {
   return unmarshal(data, v, '[')
}

func UnmarshalObject(data []byte, v interface{}) error {
   return unmarshal(data, v, '{')
}

func unmarshal(data []byte, v interface{}, c byte) error {
   for {
      ind := bytes.IndexByte(data, c)
      if ind == -1 {
         return fmt.Errorf("%q not found", c)
      }
      data = data[ind:]
      dec := json.NewDecoder(bytes.NewReader(data))
      _, err := dec.Token()
      if err == nil {
         for {
            _, err := dec.Token()
            if err != nil {
               off := dec.InputOffset()
               err := json.Unmarshal(data[:off], v)
               if err == nil {
                  return nil
               }
               data = data[off:]
               break
            }
         }
      }
   }
}
