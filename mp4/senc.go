package mp4

import (
	"fmt"
	"io"
)

func (s *sencBox) parseReadBox() error {
	if !s.readButNotParsed {
		return fmt.Errorf("senc box already parsed")
	}
	sr := newFixedSliceReader(s.rawData)
	nrBytesLeft := uint32(sr.NrRemainingBytes())
	if s.flags&useSubSampleEncryption == 0 {
		// No subsamples
		perSampleIVSize := byte(nrBytesLeft / s.SampleCount)
		s.IVs = make([]initializationVector, 0, s.SampleCount)
		switch perSampleIVSize {
		case 0:
			// Nothing to do
		case 8:
			for i := 0; i < int(s.SampleCount); i++ {
				s.IVs = append(s.IVs, sr.ReadBytes(8))
			}
		case 16:
			for i := 0; i < int(s.SampleCount); i++ {
				s.IVs = append(s.IVs, sr.ReadBytes(16))
			}
		default:
			return fmt.Errorf("strange derived PerSampleIVSize: %d", perSampleIVSize)
		}
		s.readButNotParsed = false
		return nil
	}
	startPos := sr.GetPos()
	ok := false
	for perSampleIVSize := byte(0); perSampleIVSize <= 16; perSampleIVSize += 8 {
		sr.SetPos(startPos)
		ok = s.parseAndFillSamples(sr, perSampleIVSize)
		if ok {
			break
		}
	}
	if !ok {
		return fmt.Errorf("could not decode senc")
	}
	s.readButNotParsed = false
	return nil
}

// UseSubSampleEncryption - flag for subsample encryption
const useSubSampleEncryption = 0x2

type subSamplePattern struct {
	bytesOfClearData     uint16
	BytesOfProtectedData uint32
}

// InitializationVector (8 or 16 bytes)
type initializationVector []byte

type sencBox struct {
	version          byte
	readButNotParsed bool
	perSampleIVSize  byte
	flags            uint32
	SampleCount      uint32
	startPos         uint64
	rawData          []byte
	IVs              []initializationVector // 8 or 16 bytes if present
	SubSamples       [][]subSamplePattern
}

// DecodeSencSR - box-specific decode
func decodeSencSR(hdr *boxHeader, startPos uint64, sr sliceReader) (box, error) {
	versionAndFlags := sr.ReadUint32()
	sampleCount := sr.ReadUint32()
	senc := sencBox{
		version:          byte(versionAndFlags >> 24),
		rawData:          sr.ReadBytes(hdr.payloadLen() - 8), // After the first 8 bytes of box content
		flags:            versionAndFlags & flagsMask,
		startPos:         startPos,
		SampleCount:      sampleCount,
		readButNotParsed: true,
	}

	if senc.SampleCount == 0 || len(senc.rawData) == 0 {
		senc.readButNotParsed = false
		return &senc, sr.AccError()
	}
	return &senc, sr.AccError()
}

// parseAndFillSamples - parse and fill senc samples given perSampleIVSize
func (s *sencBox) parseAndFillSamples(sr sliceReader, perSampleIVSize byte) (ok bool) {
	ok = true
	s.SubSamples = make([][]subSamplePattern, s.SampleCount)
	for i := 0; i < int(s.SampleCount); i++ {
		if perSampleIVSize > 0 {
			if sr.NrRemainingBytes() < int(perSampleIVSize) {
				ok = false
				break
			}
			s.IVs = append(s.IVs, sr.ReadBytes(int(perSampleIVSize)))
		}
		if sr.NrRemainingBytes() < 2 {
			ok = false
			break
		}
		subsampleCount := int(sr.ReadUint16())
		if sr.NrRemainingBytes() < subsampleCount*6 {
			ok = false
			break
		}
		s.SubSamples[i] = make([]subSamplePattern, subsampleCount)
		for j := 0; j < subsampleCount; j++ {
			s.SubSamples[i][j].bytesOfClearData = sr.ReadUint16()
			s.SubSamples[i][j].BytesOfProtectedData = sr.ReadUint32()
		}
	}
	if !ok || sr.NrRemainingBytes() != 0 {
		// Cleanup the IVs and SubSamples which may have been partially set
		s.IVs = nil
		s.SubSamples = nil
		ok = false
	}
	s.perSampleIVSize = byte(perSampleIVSize)
	return ok
}

// Type - box-specific type
func (s *sencBox) getType() string {
	return "senc"
}

//setSubSamplesUsedFlag - set flag if subsamples are used
func (s *sencBox) setSubSamplesUsedFlag() {
	for _, subSamples := range s.SubSamples {
		if len(subSamples) > 0 {
			s.flags |= useSubSampleEncryption
			break
		}
	}
}

// Size - box-specific type
func (s *sencBox) size() uint64 {
	if s.readButNotParsed {
		return boxHeaderSize + 8 + uint64(len(s.rawData)) // read 8 bytes after header
	}
	totalSize := boxHeaderSize + 8
	perSampleIVSize := s.getPerSampleIVSize()
	for i := 0; i < int(s.SampleCount); i++ {
		totalSize += perSampleIVSize
		if s.flags&useSubSampleEncryption != 0 {
			totalSize += 2 + 6*len(s.SubSamples[i])
		}
	}
	return uint64(totalSize)
}

func (s *sencBox) encode(w io.Writer) error {
	// First check if subsamplencryption is to be used since it influences the box size
	s.setSubSamplesUsedFlag()
	sw := newFixedSliceWriter(int(s.size()))
	err := s.encodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (s *sencBox) encodeSW(sw sliceWriter) error {
	s.setSubSamplesUsedFlag()
	encodeHeaderSW(s, sw)
	versionAndFlags := (uint32(s.version) << 24) + s.flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(s.SampleCount)
	perSampleIVSize := s.getPerSampleIVSize()
	for i := 0; i < int(s.SampleCount); i++ {
		if perSampleIVSize > 0 {
			sw.WriteBytes(s.IVs[i])
		}
		if s.flags&useSubSampleEncryption != 0 {
			sw.WriteUint16(uint16(len(s.SubSamples[i])))
			for _, subSample := range s.SubSamples[i] {
				sw.WriteUint16(subSample.bytesOfClearData)
				sw.WriteUint32(subSample.BytesOfProtectedData)
			}
		}
	}
	return sw.AccError()
}

// GetPerSampleIVSize - return perSampleIVSize
func (s *sencBox) getPerSampleIVSize() int {
	return int(s.perSampleIVSize)
}
