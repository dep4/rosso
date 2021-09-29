package json

import (
   "bytes"
   "encoding/json"
)

func Unmarshal(data []byte, v interface{}) error {
   for len(data) > 0 {
      read := bytes.NewReader(data)
      dec := json.NewDecoder(read)
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
      data = data[1:]
   }
   return NotFound{}
}

type NotFound struct{}

func (NotFound) Error() string {
   return "not found"
}
