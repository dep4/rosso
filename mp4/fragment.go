package mp4

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"sort"
)

func (f *fragment) getFullSamples() ([]fullSample, error) {
	moof := f.Moof
	mdat := f.mdat
	traf := moof.Traf // The first one
	tfhd := traf.Tfhd
	baseTime := traf.Tfdt.baseMediaDecodeTime
	moofStartPos := moof.startPos
	var samples []fullSample
	for _, trun := range traf.Truns {
		totalDur := trun.addSampleDefaultValues(tfhd)
		var baseOffset uint64
		if tfhd.HasBaseDataOffset() {
			baseOffset = tfhd.BaseDataOffset
		} else if tfhd.DefaultBaseIfMoof() {
			baseOffset = moofStartPos
		}
		if trun.hasDataOffset() {
			baseOffset = uint64(int64(trun.DataOffset) + int64(baseOffset))
		}
		mdatDataLength := uint64(len(mdat.Data)) // len should be fine for 64-bit
		offsetInMdat := baseOffset - mdat.payloadAbsoluteOffset()
		if offsetInMdat > mdatDataLength {
			return nil, errors.New("offset in mdata beyond size")
		}
		samples = append(samples, trun.getFullSamples(uint32(offsetInMdat), baseTime, mdat)...)
		baseTime += totalDur // Next trun start after this
	}
	return samples, nil
}

// DecryptBytesCTR - decrypt or encrypt sample using CTR mode, provided key, iv and sumsamplePattern
func decryptBytesCTR(Data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	inBuf := bytes.NewBuffer(Data)
	outBuf := bytes.Buffer{}

	writer := cipher.StreamWriter{S: stream, W: &outBuf}
	_, err = io.Copy(writer, inBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}

func decryptSampleCenc(sample []byte, key []byte, iv []byte, subSamplePatterns []subSamplePattern) ([]byte, error) {
	decSample := make([]byte, 0, len(sample))
	if len(subSamplePatterns) != 0 {
		var pos uint32 = 0
		for j := 0; j < len(subSamplePatterns); j++ {
			ss := subSamplePatterns[j]
			nrClear := uint32(ss.bytesOfClearData)
			nrEnc := ss.BytesOfProtectedData
			decSample = append(decSample, sample[pos:pos+nrClear]...)
			pos += nrClear
			cryptOut, err := decryptBytesCTR(sample[pos:pos+nrEnc], key, iv)
			if err != nil {
				return nil, err
			}
			decSample = append(decSample, cryptOut...)
			pos += nrEnc
		}
	} else {
		cryptOut, err := decryptBytesCTR(sample, key, iv)
		if err != nil {
			return nil, err
		}
		decSample = append(decSample, cryptOut...)
	}
	return decSample, nil
}

type fragment struct {
	Moof        *moofBox
	mdat        *mdatBox
	children    []box       // All top-level boxes in order
	encOptimize encOptimize // Bit field with optimizations being done at encoding
}

// AddChild - Add a top-level box to Fragment
func (f *fragment) addChild(b box) {
	switch b.getType() {
	case "moof":
		f.Moof = b.(*moofBox)
	case "mdat":
		f.mdat = b.(*mdatBox)
	}
	f.children = append(f.children, b)
}

func (f *fragment) encode(w io.Writer) error {
	if f.Moof == nil {
		return fmt.Errorf("moof not set in fragment")
	}
	if f.mdat == nil {
		return fmt.Errorf("mdat not set in fragment")
	}
	f.setTrunDataOffsets()
	for _, b := range f.children {
		err := b.encode(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetTrunDataOffsets - set DataOffset in trun depending on size and writeOrder
func (f *fragment) setTrunDataOffsets() {
	var truns []*trunBox
	for _, traf := range f.Moof.Trafs {
		truns = append(truns, traf.Truns...)
	}
	sort.Slice(truns, func(i, j int) bool {
		return truns[i].writeOrderNr < truns[j].writeOrderNr
	})
	dataOffset := f.Moof.size() + f.mdat.headerSize()
	for _, trun := range truns {
		trun.DataOffset = int32(dataOffset)
		dataOffset += trun.sizeOfData()
	}
}
