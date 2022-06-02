package mp4

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// FixedSliceReader - read integers and other data from a fixed slice.
// Accumulates error, and the first error can be retrieved.
// If err != nil, 0 or empty string is returned
type fixedSliceReader struct {
	err   error
	slice []byte
	pos   int
	length   int
}

// bits.NewFixedSliceReader - create a new slice reader reading from data
func newFixedSliceReader(data []byte) *fixedSliceReader {
	return &fixedSliceReader{
		slice: data,
		pos:   0,
		length:   len(data),
		err:   nil,
	}
}

// AccError - get accumulated error after read operations
func (s *fixedSliceReader) AccError() error {
	return s.err
}

// ReadUint8 - read uint8 from slice
func (s *fixedSliceReader) ReadUint8() byte {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-1 {
		s.err = errSliceRead
		return 0
	}
	res := s.slice[s.pos]
	s.pos++
	return res
}

// ReadUint16 - read uint16 from slice
func (s *fixedSliceReader) ReadUint16() uint16 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-2 {
		s.err = errSliceRead
		return 0
	}
	res := binary.BigEndian.Uint16(s.slice[s.pos : s.pos+2])
	s.pos += 2
	return res
}

// ReadInt16 - read int16 from slice
func (s *fixedSliceReader) ReadInt16() int16 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-2 {
		s.err = errSliceRead
		return 0
	}
	res := binary.BigEndian.Uint16(s.slice[s.pos : s.pos+2])
	s.pos += 2
	return int16(res)
}

// ReadUint24 - read uint24 from slice
func (s *fixedSliceReader) ReadUint24() uint32 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-3 {
		s.err = errSliceRead
		return 0
	}
	res := uint32(binary.BigEndian.Uint16(s.slice[s.pos : s.pos+2]))
	res |= res<<16 | uint32(s.slice[s.pos+2])
	s.pos += 3
	return res
}

// ReadUint32 - read uint32 from slice
func (s *fixedSliceReader) ReadUint32() uint32 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-4 {
		s.err = errSliceRead
		return 0
	}
	res := binary.BigEndian.Uint32(s.slice[s.pos : s.pos+4])
	s.pos += 4
	return res
}

// ReadInt32 - read int32 from slice
func (s *fixedSliceReader) ReadInt32() int32 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-4 {
		s.err = errSliceRead
		return 0
	}
	res := binary.BigEndian.Uint32(s.slice[s.pos : s.pos+4])
	s.pos += 4
	return int32(res)
}

// ReadUint64 - read uint64 from slice
func (s *fixedSliceReader) ReadUint64() uint64 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-8 {
		s.err = errSliceRead
		return 0
	}
	res := binary.BigEndian.Uint64(s.slice[s.pos : s.pos+8])
	s.pos += 8
	return res
}

// ReadInt64 - read int64 from slice
func (s *fixedSliceReader) ReadInt64() int64 {
	if s.err != nil {
		return 0
	}
	if s.pos > s.length-8 {
		s.err = errSliceRead
		return 0
	}
	res := binary.BigEndian.Uint64(s.slice[s.pos : s.pos+8])
	s.pos += 8
	return int64(res)
}

func (s *fixedSliceReader) ReadFixedLengthString(n int) string {
	if s.err != nil {
		return ""
	}
	if s.pos > s.length-n {
		s.err = errSliceRead
		return ""
	}
	res := string(s.slice[s.pos : s.pos+n])
	s.pos += n
	return res
}

// ReadZeroTerminatedString - read string until zero byte but at most maxLen
// Set err and return empty string if no zero byte found
func (s *fixedSliceReader) ReadZeroTerminatedString(maxLen int) string {
	if s.err != nil {
		return ""
	}
	startPos := s.pos
	maxPos := startPos + maxLen
	for {
		if s.pos >= maxPos {
			s.err = errors.New("did not find terminating zero")
			return ""
		}
		c := s.slice[s.pos]
		if c == 0 {
			str := string(s.slice[startPos:s.pos])
			s.pos++ // Next position to read
			return str
		}
		s.pos++
	}
}

// ReadBytes - read a slice of n bytes
// Return empty slice if n bytes not available
func (s *fixedSliceReader) ReadBytes(n int) []byte {
	if s.err != nil {
		return []byte{}
	}
	if s.pos > s.length-n {
		s.err = errSliceRead
		return []byte{}
	}
	res := s.slice[s.pos : s.pos+n]
	s.pos += n
	return res
}

// RemainingBytes - return remaining bytes of this slice
func (s *fixedSliceReader) RemainingBytes() []byte {
	if s.err != nil {
		return []byte{}
	}
	res := s.slice[s.pos:]
	s.pos = s.length
	return res
}

// NrRemainingBytes - return number of bytes remaining
func (s *fixedSliceReader) NrRemainingBytes() int {
	if s.err != nil {
		return 0
	}
	return s.length - s.GetPos()
}

// SkipBytes - skip passed n bytes
func (s *fixedSliceReader) SkipBytes(n int) {
	if s.err != nil {
		return
	}
	if s.pos+n > s.length {
		s.err = fmt.Errorf("attempt to skip bytes to pos %d beyond slice len %d", s.pos+n, s.length)
		return
	}
	s.pos += n
}

// SetPos - set read position is slice
func (s *fixedSliceReader) SetPos(pos int) {
	if pos > s.length {
		s.err = fmt.Errorf("attempt to set pos %d beyond slice len %d", pos, s.length)
		return
	}
	s.pos = pos
}

// GetPos - get read position is slice
func (s *fixedSliceReader) GetPos() int {
	return s.pos
}
