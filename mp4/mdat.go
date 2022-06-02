package mp4

import (
	"encoding/binary"
	"io"
)

func (m *mdatBox) encode(w io.Writer) error {
	headerSize := boxHeaderSize
	if m.largeSize {
		headerSize += 8
	}
	buf := make([]byte, headerSize)
	if m.largeSize {
		binary.BigEndian.PutUint32(buf, 1) // signals large size
		strtobuf(buf[4:], "mdat", 4)
		binary.BigEndian.PutUint64(buf[8:], m.size())
	} else {
		binary.BigEndian.PutUint32(buf, uint32(m.size()))
		strtobuf(buf[4:], "mdat", 4)
	}
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	if len(m.dataParts) > 0 {
		for _, dp := range m.dataParts {
			_, err = w.Write(dp)
			if err != nil {
				return err
			}
		}
	} else {
		_, err = w.Write(m.Data)
	}
	return err
}

func decodeMdat(head *boxHeader, startPos uint64, r io.Reader) (box, error) {
	data, err := head.readBoxBody(r)
	if err != nil {
		return nil, err
	}
	largeSize := head.length > boxHeaderSize
	return &mdatBox{startPos, data, nil, 0, largeSize}, nil
}

// TfdtBox - Track Fragment Decode Time (tfdt)
//
// Contained in : Track Fragment box (traf)
type tfdtBox struct {
	version             byte
	flags               uint32
	baseMediaDecodeTime uint64
}

// DecodeTfdtSR - box-specific decode
func decodeTfdtSR(hdr *boxHeader, startPos uint64, sr sliceReader) (box, error) {
	versionAndFlags := sr.ReadUint32()
	version := byte(versionAndFlags >> 24)
	var baseMediaDecodeTime uint64
	if version == 0 {
		baseMediaDecodeTime = uint64(sr.ReadUint32())
	} else {
		baseMediaDecodeTime = sr.ReadUint64()
	}

	b := tfdtBox{
		version:             version,
		flags:               versionAndFlags & flagsMask,
		baseMediaDecodeTime: baseMediaDecodeTime,
	}
	return &b, sr.AccError()
}

// Type - return box type
func (t *tfdtBox) getType() string {
	return "tfdt"
}

// Size - return calculated size
func (t *tfdtBox) size() uint64 {
	return uint64(boxHeaderSize + 8 + 4*int(t.version))
}

// Encode - write box to w
func (t *tfdtBox) encode(w io.Writer) error {
	sw := newFixedSliceWriter(int(t.size()))
	err := t.encodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (t *tfdtBox) encodeSW(sw sliceWriter) error {
	encodeHeaderSW(t, sw)
	versionAndFlags := (uint32(t.version) << 24) + t.flags
	sw.WriteUint32(versionAndFlags)
	if t.version == 0 {
		sw.WriteUint32(uint32(t.baseMediaDecodeTime))
	} else {
		sw.WriteUint64(t.baseMediaDecodeTime)
	}
	return sw.AccError()
}

// MdatBox - Media Data Box (mdat)
// The mdat box contains media chunks/samples.
// DataParts is to be able to gather output data without
// new allocations
type mdatBox struct {
	startPos     uint64
	Data         []byte
	dataParts    [][]byte
	lazyDataSize uint64
	largeSize    bool
}

const maxNormalPayloadSize = (1 << 32) - 1 - 8

// Type - return box type
func (m *mdatBox) getType() string {
	return "mdat"
}

// Size - return calculated size, depending on largeSize set or not
func (m *mdatBox) size() uint64 {
	dataSize := m.dataLength()

	if m.lazyDataSize > 0 {
		dataSize = m.lazyDataSize
	}
	if dataSize > maxNormalPayloadSize {
		m.largeSize = true
	}
	size := boxHeaderSize + dataSize
	if m.largeSize {
		size += 8
	}
	return size
}

// DataLength - length of data stored in box either as one or multiple parts
func (m *mdatBox) dataLength() uint64 {
	dataLength := len(m.Data)
	if len(m.dataParts) > 0 {
		dataLength = 0
		for i := range m.dataParts {
			dataLength += len(m.dataParts[i])
		}
	}
	return uint64(dataLength)
}

// HeaderSize - 8 or 16 (bytes) depending o whether largeSize is used
func (m *mdatBox) headerSize() uint64 {
	hSize := boxHeaderSize
	if m.largeSize {
		hSize += largeSizeLen
	}
	return uint64(hSize)
}

// PayloadAbsoluteOffset - position of mdat payload start (works after header)
func (m *mdatBox) payloadAbsoluteOffset() uint64 {
	return m.startPos + m.headerSize()
}
