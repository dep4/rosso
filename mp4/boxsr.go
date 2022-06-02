package mp4

import (
	"fmt"
	"io"
)

func decodeBoxSR(startPos uint64, sr sliceReader) (box, error) {
   h := new(boxHeader)
   // first
   h.size = uint64(sr.ReadUint32())
   // second
   h.name = sr.ReadFixedLengthString(4)
   h.length = boxHeaderSize
   var (
      b box
      err error
   )
   switch h.name {
   case "senc":
      b, err = decodeSencSR(h, startPos, sr)
   case "tfdt":
      b, err = decodeTfdtSR(h, startPos, sr)
   case "tfhd":
      b, err = decodeTfhdSR(h, startPos, sr)
   case "traf":
      b, err = decodeTrafSR(h, startPos, sr)
   case "trun":
      b, err = decodeTrunSR(h, startPos, sr)
   default:
      b, err = decodeUnknownSR(h, startPos, sr)
   }
   if err != nil {
      return nil, fmt.Errorf("decode %s: %w", h.name, err)
   }
   return b, nil
}

// DecodeUnknown - decode an unknown box
func decodeUnknownSR(hdr *boxHeader, startPos uint64, sr sliceReader) (box, error) {
	return &unknownBox{hdr.name, hdr.size, sr.ReadBytes(hdr.payloadLen())}, sr.AccError()
}

func decodeUnknown(hdr *boxHeader, startPos uint64, r io.Reader) (box, error) {
	data, err := hdr.readBoxBody(r)
	if err != nil {
		return nil, err
	}
	sr := newFixedSliceReader(data)
	return decodeUnknownSR(hdr, startPos, sr)
}

// UnknownBox - box that we don't know how to parse
type unknownBox struct {
	name       string
	length     uint64
	notDecoded []byte
}

// Type - return box type
func (b *unknownBox) getType() string {
	return b.name
}

// Size - return calculated size
func (b *unknownBox) size() uint64 {
	return b.length
}

func (b *unknownBox) encode(w io.Writer) error {
	sw := newFixedSliceWriter(int(b.size()))
	err := b.encodeSW(sw)
	if err != nil {
		return err
	}
	if _, err := w.Write(sw.Bytes()); err != nil {
		return err
	}
	return nil
}

func (b *unknownBox) encodeSW(sw sliceWriter) error {
	encodeHeaderSW(b, sw)
	sw.WriteBytes(b.notDecoded)
	return sw.AccError()
}

// Sample - sample as used in trun box (mdhd timescale)
type sample struct {
	flags                 uint32 // interpreted as SampleFlags
	dur                   uint32 // Sample duration in mdhd timescale
	size                  uint32 // Size of sample data
	compositionTimeOffset int32  // Signed composition time offset
}

type fullSample struct {
	sample
	decodeTime uint64 // Absolute decode time (offset + accumulated sample Dur)
	Data       []byte // Sample data
}

type mediaSegment struct {
	encOptimize encOptimize
	Fragments   []*fragment
}

func (s *mediaSegment) encode(w io.Writer) error {
   for _, f := range s.Fragments {
      err := f.encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}
