package protobuf

import (
   "bufio"
   "bytes"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strconv"
)

func ConsumeGroup(num Number, b []byte) (v []byte, n int) {
   n = ConsumeFieldValue(num, StartGroupType, b)
   if n < 0 {
      return nil, n // forward error code
   }
   b = b[:n]
   for len(b) > 0 && b[len(b)-1]&0x7f == 0 {
      b = b[:len(b)-1]
   }
   b = b[:len(b)-SizeTag(num)]
   return b, n
}

func ConsumeFieldValue(num Number, typ Type, b []byte) (n int) {
   return consumeFieldValueD(num, typ, b, DefaultRecursionLimit)
}

func consumeFieldValueD(num Number, typ Type, b []byte, depth int) (n int) {
	switch typ {
	case VarintType:
		_, n = ConsumeVarint(b)
		return n
	case Fixed32Type:
		_, n = ConsumeFixed32(b)
		return n
	case Fixed64Type:
		_, n = ConsumeFixed64(b)
		return n
	case BytesType:
		_, n = ConsumeBytes(b)
		return n
	case StartGroupType:
		if depth < 0 {
			return errCodeRecursionDepth
		}
		n0 := len(b)
		for {
			num2, typ2, n := ConsumeTag(b)
			if n < 0 {
				return n // forward error code
			}
			b = b[n:]
			if typ2 == EndGroupType {
				if num != num2 {
					return errCodeEndGroup
				}
				return n0 - len(b)
			}

			n = consumeFieldValueD(num2, typ2, b, depth-1)
			if n < 0 {
				return n // forward error code
			}
			b = b[n:]
		}
	case EndGroupType:
		return errCodeEndGroup
	default:
		return errCodeReserved
	}
}
