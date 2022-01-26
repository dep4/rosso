package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) Add(num protowire.Number, name string, val Message) error {
   if m == nil {
      return nilMap{"protobuf.Message.Add"}
   }
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{value, val}
   case []Message:
      m[tag] = append(value, val)
   }
   return nil
}

func (m Message) Get(num protowire.Number, name string) Message {
   for _, name := range []string{name, messageType} {
      tag := Tag{num, name}
      val, ok := m[tag].(Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number, name string) []byte {
   for _, name := range []string{name, bytesType} {
      tag := Tag{num, name}
      val, ok := m[tag].([]byte)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetFixed64(num protowire.Number, name string) uint64 {
   for _, name := range []string{name, fixed64Type} {
      tag := Tag{num, name}
      val, ok := m[tag].(uint64)
      if ok {
         return val
      }
   }
   return 0
}

func (m Message) GetMessages(num protowire.Number, name string) []Message {
   for _, name := range []string{name, messageType} {
      tag := Tag{num, name}
      val, ok := m[tag].([]Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetString(num protowire.Number, name string) string {
   for _, name := range []string{name, stringType} {
      tag := Tag{num, name}
      val, ok := m[tag].(string)
      if ok {
         return val
      }
   }
   return ""
}

func (m Message) GetVarint(num protowire.Number, name string) uint64 {
   for _, name := range []string{name, varintType} {
      tag := Tag{num, name}
      val, ok := m[tag].(uint64)
      if ok {
         return val
      }
   }
   return 0
}

func (m Message) addString(num protowire.Number, val string) {
   tag := Tag{num, stringType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{value, val}
   case []string:
      m[tag] = append(value, val)
   }
}

func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(elem2)
   mes, err := Unmarshal(elem2)
   if err != nil {
      if binary {
         tag := Tag{num, bytesType}
         switch elem := m[tag].(type) {
         case nil:
            m[tag] = elem2
         case []byte:
            m[tag] = [][]byte{elem, elem2}
         case [][]byte:
            m[tag] = append(elem, elem2)
         }
      } else {
         m.addString(num, string(elem2))
      }
   } else {
      // In this section, if input is binary, then result could be a Message or
      // []byte. We assume for now its always a Message. If input is not binary,
      // then result could be a Message or string. Since its not possible to
      // tell Message from string, we just add both under the same Number, each
      // with its own Type.
      m.Add(num, "", mes)
      if !binary {
         m.addString(num, string(elem2))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   tag := Tag{num, fixed64Type}
   switch elem := m[tag].(type) {
   case nil:
      m[tag] = elem2
   case uint64:
      m[tag] = []uint64{elem, elem2}
   case []uint64:
      m[tag] = append(elem, elem2)
   }
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   elem2, eLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   tag := Tag{num, varintType}
   switch elem := m[tag].(type) {
   case nil:
      m[tag] = elem2
   case uint64:
      m[tag] = []uint64{elem, elem2}
   case []uint64:
      m[tag] = append(elem, elem2)
   }
   return nil
}
