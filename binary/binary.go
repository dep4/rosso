// Simple translation between numbers and byte sequences.
//
// Why not use existing packages? These are not big-endian:
//
//  encoding/binary#Uvarint
//  encoding/binary#Varint
//
// These read a fixed number of bytes:
//
//  encoding/binary#ByteOrder.Uint16
//  encoding/binary#ByteOrder.Uint32
//  encoding/binary#ByteOrder.Uint64
//  x/crypto/cryptobyte#String.ReadUint8
//  x/crypto/cryptobyte#String.ReadUint16
//  x/crypto/cryptobyte#String.ReadUint32
package binary

import (
   "encoding/binary"
)

var (
   Uint16 = binary.BigEndian.Uint16
   Uint32 = binary.BigEndian.Uint32
   Uint64 = binary.BigEndian.Uint64
)

func UintN(buf []byte, n int) uint64 {
   var length uint64
   for _, b := range buf[:n] {
      length <<= 8
      length |= uint64(b)
   }
   return length
}

func IntN(buf []byte, n int) int64 {
   var length int64
   for _, b := range buf[:n] {
      length <<= 8
      length |= int64(b)
   }
   return length
}
