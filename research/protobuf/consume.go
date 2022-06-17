package protobuf

import (
   "bufio"
   "bytes"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func consumeBytes(buf *bufio.Reader) (Bytes, error) {
   var val Bytes
   n, err := binary.ReadUvarint(buf)
   if err != nil {
      return val, err
   }
   var limit io.LimitedReader
   limit.N = int64(n)
   limit.R = buf
   val.Raw, err = io.ReadAll(&limit)
   if err != nil {
      return val, err
   }
   val.Message, _ = Decode(bytes.NewReader(val.Raw))
   return val, nil
}

func consumeFixed32(buf io.Reader) (Fixed32, error) {
   var val Fixed32
   err := binary.Read(buf, binary.LittleEndian, &val)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeFixed64(buf io.Reader) (Fixed64, error) {
   var val Fixed64
   err := binary.Read(buf, binary.LittleEndian, &val)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeVarint(buf io.ByteReader) (Varint, error) {
   val, err := binary.ReadUvarint(buf)
   if err != nil {
      return 0, err
   }
   return Varint(val), nil
}

func consumeTag(buf io.ByteReader) (Number, protowire.Type, error) {
   tag, err := binary.ReadUvarint(buf)
   if err != nil {
      return 0, 0, err
   }
   num, typ := protowire.DecodeTag(tag)
   if num < protowire.MinValidNumber {
      return 0, 0, errors.New("invalid field number")
   }
   return num, typ, nil
}
