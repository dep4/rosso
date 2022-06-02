package mp4

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type sliceReader interface {
	AccError() error
	GetPos() int
	NrRemainingBytes() int
	ReadBytes(n int) []byte
	ReadFixedLengthString(n int) string
	ReadInt32() int32
	ReadInt64() int64
	ReadUint16() uint16
	ReadUint32() uint32
	ReadUint64() uint64
	ReadUint8() byte
	ReadZeroTerminatedString(maxLen int) string
	SkipBytes(n int)
}

type sliceWriter interface {
	AccError() error
	WriteBits(bits uint, n int)
	WriteBytes(byteSlice []byte)
	WriteInt16(n int16)
	WriteInt32(n int32)
	WriteInt64(n int64)
	WriteString(s string, addZeroEnd bool)
	WriteUint16(n uint16)
	WriteUint32(n uint32)
	WriteUint64(n uint64)
	WriteUint8(n byte)
}

// SliceReader errors
var errSliceRead = fmt.Errorf("read too far in SliceReader")

var errSliceWriter = errors.New("overflow in SliceWriter")

// mask - n-bit binary mask
func mask(n int) uint {
	return (1 << uint(n)) - 1
}

// FixedSliceWriter - write numbers to a fixed []byte slice
type fixedSliceWriter struct {
	accError error
	buf      []byte
	off      int
	n        int  // current number of bits
	v        uint // current accumulated value for bits
}

func (sw *fixedSliceWriter) WriteBits(bits uint, n int) {
	if sw.accError != nil {
		return
	}
	sw.v <<= uint(n)
	sw.v |= bits & mask(n)
	sw.n += n
	for sw.n >= 8 {
		b := byte((sw.v >> (uint(sw.n) - 8)) & mask(8))
		sw.WriteUint8(b)
		sw.n -= 8
	}
	sw.v &= mask(8)
}

// NewSliceWriter - create slice writer with fixed size.
func newFixedSliceWriter(size int) *fixedSliceWriter {
	return &fixedSliceWriter{
		buf:      make([]byte, size),
		off:      0,
		n:        0,
		v:        0,
		accError: nil,
	}
}

// Bytes - return buf up to what's written
func (sw *fixedSliceWriter) Bytes() []byte {
	return sw.buf[:sw.off]
}

// Offset - offset for writing in FixedSliceWriter buffer
func (sw *fixedSliceWriter) Offset() int {
	return sw.off
}

// AccError - return accumulated erro
func (sw *fixedSliceWriter) AccError() error {
	return sw.accError
}

// WriteUint8 - write byte to slice
func (sw *fixedSliceWriter) WriteUint8(n byte) {
	if sw.off+1 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	sw.buf[sw.off] = n
	sw.off++
}

// WriteUint16 - write uint16 to slice
func (sw *fixedSliceWriter) WriteUint16(n uint16) {
	if sw.off+2 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	binary.BigEndian.PutUint16(sw.buf[sw.off:], n)
	sw.off += 2
}

// WriteInt16 - write int16 to slice
func (sw *fixedSliceWriter) WriteInt16(n int16) {
	if sw.off+2 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	binary.BigEndian.PutUint16(sw.buf[sw.off:], uint16(n))
	sw.off += 2
}

// WriteUint32 - write uint32 to slice
func (sw *fixedSliceWriter) WriteUint32(n uint32) {
	if sw.off+4 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	binary.BigEndian.PutUint32(sw.buf[sw.off:], n)
	sw.off += 4
}

// WriteInt32 - write int32 to slice
func (sw *fixedSliceWriter) WriteInt32(n int32) {
	if sw.off+4 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	binary.BigEndian.PutUint32(sw.buf[sw.off:], uint32(n))
	sw.off += 4
}

// WriteUint64 - write uint64 to slice
func (sw *fixedSliceWriter) WriteUint64(n uint64) {
	if sw.off+8 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	binary.BigEndian.PutUint64(sw.buf[sw.off:], n)
	sw.off += 8
}

// WriteInt64 - write int64 to slice
func (sw *fixedSliceWriter) WriteInt64(n int64) {
	if sw.off+8 > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	binary.BigEndian.PutUint64(sw.buf[sw.off:], uint64(n))
	sw.off += 8
}

// WriteString - write string to slice with or without zero end
func (sw *fixedSliceWriter) WriteString(s string, addZeroEnd bool) {
	nrNew := len(s)
	if addZeroEnd {
		nrNew++
	}
	if sw.off+nrNew > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	copy(sw.buf[sw.off:sw.off+len(s)], s)
	sw.off += len(s)
	if addZeroEnd {
		sw.buf[sw.off] = 0
		sw.off++
	}
}

// WriteZeroBytes - write n byte of zeroes
func (sw *fixedSliceWriter) WriteZeroBytes(n int) {
	if sw.off+n > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	for i := 0; i < n; i++ {
		sw.buf[sw.off] = 0
		sw.off++
	}
}

// WriteBytes - write []byte
func (sw *fixedSliceWriter) WriteBytes(byteSlice []byte) {
	if sw.off+len(byteSlice) > len(sw.buf) {
		sw.accError = errSliceWriter
		return
	}
	copy(sw.buf[sw.off:sw.off+len(byteSlice)], byteSlice)
	sw.off += len(byteSlice)
}
