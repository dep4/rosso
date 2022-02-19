package protobuf

type message map[float64]interface{}

func set() {
   mes := make(message)
   mes.add(4, make(message))
   mes.get(4).set(1, make(message))
   mes.get(4).get(1).setUint64(10, 29)
}
