package protobuf

import (
   "bufio"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func consumeBytes(buf *bufio.Reader) ([]byte, error) {
   n, err := binary.ReadUvarint(buf)
   if err != nil {
      return nil, err
   }
   var limit io.LimitedReader
   limit.N = int64(n)
   limit.R = buf
   return io.ReadAll(&limit)
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
