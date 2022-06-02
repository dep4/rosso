package mp4

import (
	"io"
)

// TrunBox - Track Fragment Run Box (trun)
//
// Contained in :  Track Fragmnet Box (traf)
//
type trunBox struct {
	version          byte
	flags            uint32
	sampleCount      uint32
	DataOffset       int32
	firstSampleFlags uint32 // interpreted same way as SampleFlags
	Samples          []sample
	writeOrderNr     uint32 // Used for multi trun offsets
}

const trunDataOffsetPresentFlag uint32 = 0x01
const trunFirstSampleFlagsPresentFlag uint32 = 0x04
const trunSampleDurationPresentFlag uint32 = 0x100
const trunSampleSizePresentFlag uint32 = 0x200
const trunSampleFlagsPresentFlag uint32 = 0x400
const trunSampleCompositionTimeOffsetPresentFlag uint32 = 0x800

// DecodeTrun - box-specific decode
func decodeTrunSR(hdr *boxHeader, startPos uint64, sr sliceReader) (box, error) {
	versionAndFlags := sr.ReadUint32()
	sampleCount := sr.ReadUint32()
	t := &trunBox{
		version:     byte(versionAndFlags >> 24),
		flags:       versionAndFlags & flagsMask,
		sampleCount: sampleCount,
		Samples:     make([]sample, sampleCount),
	}

	if t.hasDataOffset() {
		t.DataOffset = sr.ReadInt32()
	}

	if t.hasFirstSampleFlags() {
		t.firstSampleFlags = sr.ReadUint32()
	}

	var i uint32
	for i = 0; i < t.sampleCount; i++ {
		var dur, size, flags uint32
		var cto int32
		if t.hasSampleDuration() {
			dur = sr.ReadUint32()
		}
		if t.hasSampleSize() {
			size = sr.ReadUint32()
		}
		if t.hasSampleFlags() {
			flags = sr.ReadUint32()
		} else if t.hasFirstSampleFlags() && i == 0 {
			flags = t.firstSampleFlags
		}
		if t.hasSampleCompositionTimeOffset() {
			cto = sr.ReadInt32()
		}
		t.Samples[i] = sample{flags, dur, size, cto}
	}

	return t, sr.AccError()
}

// AddSampleDefaultValues - add values from tfhd and trex boxes if needed
// Return total duration
func (t *trunBox) addSampleDefaultValues(tfhd *tfhdBox) (totalDur uint64) {
	var defaultSampleDuration uint32
	var defaultSampleSize uint32
	var defaultSampleFlags uint32
	if tfhd.HasDefaultSampleDuration() {
		defaultSampleDuration = tfhd.DefaultSampleDuration
	}
	if tfhd.HasDefaultSampleSize() {
		defaultSampleSize = tfhd.DefaultSampleSize
	}
	if tfhd.HasDefaultSampleFlags() {
		defaultSampleFlags = tfhd.DefaultSampleFlags
	}
	var i uint32
	totalDur = 0
	for i = 0; i < t.sampleCount; i++ {
		if !t.hasSampleDuration() {
			t.Samples[i].dur = defaultSampleDuration
		}
		totalDur += uint64(t.Samples[i].dur)
		if !t.hasSampleSize() {
			t.Samples[i].size = defaultSampleSize
		}
		if !t.hasSampleFlags() {
			if i > 0 || !t.hasFirstSampleFlags() {
				t.Samples[i].flags = defaultSampleFlags
			}
		}
	}
	return totalDur
}

// HasDataOffset - interpreted dataOffsetPresent flag
func (t *trunBox) hasDataOffset() bool {
	return t.flags&trunDataOffsetPresentFlag != 0
}

// HasFirstSampleFlags - interpreted firstSampleFlagsPresent flag
func (t *trunBox) hasFirstSampleFlags() bool {
	return t.flags&trunFirstSampleFlagsPresentFlag != 0
}

// HasSampleDuration - interpreted sampleDurationPresent flag
func (t *trunBox) hasSampleDuration() bool {
	return t.flags&trunSampleDurationPresentFlag != 0
}

// HasSampleFlags - interpreted sampleFlagsPresent flag
func (t *trunBox) hasSampleFlags() bool {
	return t.flags&trunSampleFlagsPresentFlag != 0
}

// HasSampleSize - interpreted sampleSizePresent flag
func (t *trunBox) hasSampleSize() bool {
	return t.flags&trunSampleSizePresentFlag != 0
}

// HasSampleCompositionTimeOffset - interpreted sampleCompositionTimeOffset flag
func (t *trunBox) hasSampleCompositionTimeOffset() bool {
	return t.flags&trunSampleCompositionTimeOffsetPresentFlag != 0
}

// Type - return box type
func (t *trunBox) getType() string {
	return "trun"
}

// Size - return calculated size
func (t *trunBox) size() uint64 {
	sz := boxHeaderSize + 8 // flags + entrycCount
	if t.hasDataOffset() {
		sz += 4
	}
	if t.hasFirstSampleFlags() {
		sz += 4
	}
	bytesPerSample := 0
	if t.hasSampleDuration() {
		bytesPerSample += 4
	}
	if t.hasSampleSize() {
		bytesPerSample += 4
	}
	if t.hasSampleFlags() {
		bytesPerSample += 4
	}
	if t.hasSampleCompositionTimeOffset() {
		bytesPerSample += 4
	}
	sz += int(t.sampleCount) * bytesPerSample
	return uint64(sz)
}

// Encode - write box to w
func (t *trunBox) encode(w io.Writer) error {
	sw := newFixedSliceWriter(int(t.size()))
	err := t.encodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (t *trunBox) encodeSW(sw sliceWriter) error {
	encodeHeaderSW(t, sw)
	versionAndFlags := (uint32(t.version) << 24) + t.flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(t.sampleCount)
	if t.hasDataOffset() {
		if t.DataOffset == 0 {
			panic("trun data offset not set")
		}
		sw.WriteInt32(t.DataOffset)
	}
	if t.hasFirstSampleFlags() {
		sw.WriteUint32(t.firstSampleFlags)
	}
	var i uint32
	for i = 0; i < t.sampleCount; i++ {
		if t.hasSampleDuration() {
			sw.WriteUint32(t.Samples[i].dur)
		}
		if t.hasSampleSize() {
			sw.WriteUint32(t.Samples[i].size)
		}
		if t.hasSampleFlags() {
			sw.WriteUint32(t.Samples[i].flags)
		}
		if t.hasSampleCompositionTimeOffset() {
			sw.WriteInt32(t.Samples[i].compositionTimeOffset)
		}

	}
	return sw.AccError()
}

// GetFullSamples - get all sample data including accumulated time and binary
// media data offsetInMdat is offset in mdat data (data normally starts 8 or 16
// bytes after start of mdat box) baseDecodeTime is decodeTime in tfdt in track
// timescale (timescale in mfhd) To fill missing individual values from tfhd
// and trex defaults, call trun.AddSampleDefaultValues() before this call
func (t *trunBox) getFullSamples(offsetInMdat uint32, baseDecodeTime uint64, mdat *mdatBox) []fullSample {
	samples := make([]fullSample, 0)
	var accDur uint64 = 0
	for _, s := range t.Samples {
		dTime := baseDecodeTime + accDur

		newSample := fullSample{
			sample:     s,
			decodeTime: dTime,
			Data:       mdat.Data[offsetInMdat : offsetInMdat+s.size],
		}
		samples = append(samples, newSample)
		accDur += uint64(s.dur)
		offsetInMdat += s.size
	}
	return samples
}

// SizeOfData - size of mediasamples in bytes
func (t *trunBox) sizeOfData() (totalSize uint64) {
	for _, sample := range t.Samples {
		totalSize += uint64(sample.size)
	}
	return totalSize
}
