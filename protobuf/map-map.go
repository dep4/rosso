package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

func newMessage(m map[string]interface{}) (Message, error) {
   mes := make(Message)
   for str, val := range m {
      num, err := strconv.Atoi(str)
      if err != nil {
         return nil, err
      }
      if err := mes.set(protowire.Number(num), val); err != nil {
         return nil, err
      }
   }
   return mes, nil
}

func (m Message) set(num protowire.Number, any interface{}) error {
   switch val := any.(type) {
   case map[string]interface{}:
      mes, err := newMessage(val)
      if err != nil {
         return err
      }
      m[num] = mes
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
