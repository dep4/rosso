package protobuf

import (
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func consumeFixed64(r io.Reader) (uint64, error) {
   var v uint64
   err := binary.Read(r, binary.LittleEndian, &v)
   if err != nil {
      return 0, err
   }
   return v, nil
}

func consumeVarint(r io.ByteReader) (uint64, error) {
   return binary.ReadUvarint(r)
}

func consumeTag(r io.ByteReader) (protowire.Number, protowire.Type, error) {
   v, err := consumeVarint(r)
   if err != nil {
      return 0, 0, err
   }
   num, typ := protowire.DecodeTag(v)
   if num < protowire.MinValidNumber {
      return 0, 0, errors.New("invalid field number")
   }
   return num, typ, nil
}
