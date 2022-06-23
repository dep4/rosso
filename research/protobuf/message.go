package protobuf

type Message map[Number]Encoder

func (m Message) Add_Varint(num Number, v Varint) error {
   return add(m, num, v)
}

func (m Message) Varint(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Varint)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) Add_Fixed64(num Number, v Fixed64) error {
   return add(m, num, v)
}

func (m Message) Fixed64(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Fixed64)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) Add_Fixed32(num Number, v Fixed32) error {
   return add(m, num, v)
}

func (m Message) Add_Bytes(num Number, v Bytes) error {
   return add(m, num, v)
}

func (m Message) Add_String(num Number, v String) error {
   return add(m, num, v)
}
