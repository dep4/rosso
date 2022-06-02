package mp4

import (
	"io"
)

const baseDataOffsetPresent uint32 = 0x000001
const sampleDescriptionIndexPresent uint32 = 0x000002
const defaultSampleDurationPresent uint32 = 0x000008
const defaultSampleSizePresent uint32 = 0x000010
const defaultSampleFlagsPresent uint32 = 0x000020
const durationIsEmpty uint32 = 0x010000
const defaultBaseIsMoof uint32 = 0x020000

// TfhdBox - Track Fragment Header Box (tfhd)
//
// Contained in : Track Fragment box (traf))
type tfhdBox struct {
	version                byte
	flags                  uint32
	TrackID                uint32
	BaseDataOffset         uint64
	SampleDescriptionIndex uint32
	DefaultSampleDuration  uint32
	DefaultSampleSize      uint32
	DefaultSampleFlags     uint32
}

// DecodeTfhdSR - box-specific decode
func decodeTfhdSR(hdr *boxHeader, startPos uint64, sr sliceReader) (box, error) {
	versionAndFlags := sr.ReadUint32()
	version := byte(versionAndFlags >> 24)
	flags := versionAndFlags & flagsMask

	t := &tfhdBox{
		version: version,
		flags:   flags,
		TrackID: sr.ReadUint32(),
	}

	if t.HasBaseDataOffset() {
		t.BaseDataOffset = sr.ReadUint64()
	}
	if t.HasSampleDescriptionIndex() {
		t.SampleDescriptionIndex = sr.ReadUint32()
	}
	if t.HasDefaultSampleDuration() {
		t.DefaultSampleDuration = sr.ReadUint32()
	}
	if t.HasDefaultSampleSize() {
		t.DefaultSampleSize = sr.ReadUint32()
	}
	if t.HasDefaultSampleFlags() {
		t.DefaultSampleFlags = sr.ReadUint32()
	}

	return t, sr.AccError()
}

// HasBaseDataOffset - interpreted flags value
func (t *tfhdBox) HasBaseDataOffset() bool {
	return t.flags&baseDataOffsetPresent != 0
}

// HasSampleDescriptionIndex - interpreted flags value
func (t *tfhdBox) HasSampleDescriptionIndex() bool {
	return t.flags&sampleDescriptionIndexPresent != 0
}

// HasDefaultSampleDuration - interpreted flags value
func (t *tfhdBox) HasDefaultSampleDuration() bool {
	return t.flags&defaultSampleDurationPresent != 0
}

// HasDefaultSampleSize - interpreted flags value
func (t *tfhdBox) HasDefaultSampleSize() bool {
	return t.flags&defaultSampleSizePresent != 0
}

// HasDefaultSampleFlags - interpreted flags value
func (t *tfhdBox) HasDefaultSampleFlags() bool {
	return t.flags&defaultSampleFlagsPresent != 0
}

// DurationIsEmpty - interpreted flags value
func (t *tfhdBox) DurationIsEmpty() bool {
	return t.flags&durationIsEmpty != 0
}

// DefaultBaseIfMoof - interpreted flags value
func (t *tfhdBox) DefaultBaseIfMoof() bool {
	return t.flags&defaultBaseIsMoof != 0
}

// Type - returns box type
func (t *tfhdBox) getType() string {
	return "tfhd"
}

// Size - returns calculated size
func (t *tfhdBox) size() uint64 {
	sz := boxHeaderSize + 8
	if t.HasBaseDataOffset() {
		sz += 8
	}
	if t.HasSampleDescriptionIndex() {
		sz += 4
	}
	if t.HasDefaultSampleDuration() {
		sz += 4
	}
	if t.HasDefaultSampleSize() {
		sz += 4
	}
	if t.HasDefaultSampleFlags() {
		sz += 4
	}
	return uint64(sz)
}

// Encode - write box to w
func (t *tfhdBox) encode(w io.Writer) error {
	sw := newFixedSliceWriter(int(t.size()))
	err := t.encodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (t *tfhdBox) encodeSW(sw sliceWriter) error {
	encodeHeaderSW(t, sw)
	versionAndFlags := (uint32(t.version) << 24) + t.flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(t.TrackID)
	if t.HasBaseDataOffset() {
		sw.WriteUint64(t.BaseDataOffset)
	}
	if t.HasSampleDescriptionIndex() {
		sw.WriteUint32(t.SampleDescriptionIndex)
	}
	if t.HasDefaultSampleDuration() {
		sw.WriteUint32(t.DefaultSampleDuration)
	}
	if t.HasDefaultSampleSize() {
		sw.WriteUint32(t.DefaultSampleSize)
	}
	if t.HasDefaultSampleFlags() {
		sw.WriteUint32(t.DefaultSampleFlags)
	}
	return sw.AccError()
}
