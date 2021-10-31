// Simple translation between numbers and byte sequences.
//
// Why not use existing packages? These are not big-endian:
//
//  encoding/binary#Uvarint
//  encoding/binary#Varint
//
// These read a fixed number of bytes:
//
//  x/crypto/cryptobyte#String.ReadUint8
//  x/crypto/cryptobyte#String.ReadUint16
//  x/crypto/cryptobyte#String.ReadUint32
//
// These panic if buffer is short:
//
//  encoding/binary#ByteOrder.Uint16
//  encoding/binary#ByteOrder.Uint32
//  encoding/binary#ByteOrder.Uint64
package binary

import (
   "encoding/binary"
)

func Uint16(buf []byte) (uint16, bool) {
   if len(buf) < 2 {
      return 0, false
   }
   return binary.BigEndian.Uint16(buf), true
}

func Uint32(buf []byte) (uint32, bool) {
   if len(buf) < 4 {
      return 0, false
   }
   return binary.BigEndian.Uint32(buf), true
}

func Uint64(buf []byte) (uint64, bool) {
   if len(buf) < 8 {
      return 0, false
   }
   return binary.BigEndian.Uint64(buf), true
}

func UintN(buf []byte, n int) (uint64, bool) {
   if n < 0 || n > len(buf) {
      return 0, false
   }
   var length uint64
   for _, b := range buf[:n] {
      length <<= 8
      length |= uint64(b)
   }
   return length, true
}
