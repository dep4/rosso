package binary

func Int16(buf [2]byte) int16 {
   var length int16
   for _, b := range buf {
      length <<= 8
      length |= int16(b)
   }
   return length
}

func Uint16(buf [2]byte) uint16 {
   var length uint16
   for _, b := range buf {
      length <<= 8
      length |= uint16(b)
   }
   return length
}

func Uvarint(buf []byte) uint64 {
   var length uint64
   for _, b := range buf {
      length <<= 8
      length |= uint64(b)
   }
   return length
}

func Varint(buf []byte) int64 {
   var length int64
   for _, b := range buf {
      length <<= 8
      length |= int64(b)
   }
   return length
}
