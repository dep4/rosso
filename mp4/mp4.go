package mp4

import (
	"fmt"
	"io"
)

func DecryptMP4withCenc(r io.Reader, key []byte, w io.Writer) error {
	inMp4, err := decodeFile(r)
	if err != nil {
		return err
	}
	if !inMp4.IsFragmented {
		return fmt.Errorf("file not fragmented. Not supported")
	}
	for _, seg := range inMp4.Segments {
		for _, frag := range seg.Fragments {
			for _, traf := range frag.Moof.Trafs {
				samples, err := frag.getFullSamples()
				if err != nil {
					return err
				}
				for i := range samples {
					encSample := samples[i].Data
					var iv []byte
					if len(traf.Senc.IVs[i]) == 8 {
						iv = make([]byte, 0, 16)
						iv = append(iv, traf.Senc.IVs[i]...)
						iv = append(iv, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
					} else {
						iv = traf.Senc.IVs[i]
					}
					var subSamplePatterns []subSamplePattern
					if len(traf.Senc.SubSamples) != 0 {
						subSamplePatterns = traf.Senc.SubSamples[i]
					}
					decryptedSample, err := decryptSampleCenc(encSample, key, iv, subSamplePatterns)
					if err != nil {
						return err
					}
					_ = copy(samples[i].Data, decryptedSample)
				}
				traf.removeEncryptionBoxes()
			}
		}
		err := seg.encode(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeMoof(hdr *boxHeader, startPos uint64, r io.Reader) (box, error) {
	Data := make([]byte, hdr.payloadLen())
	_, err := io.ReadFull(r, Data)
	if err != nil {
		return nil, err
	}
	sr := newFixedSliceReader(Data)
	children, err := decodeContainerChildrenSR(hdr, startPos+8, startPos+hdr.size, sr)
	if err != nil {
		return nil, err
	}
	m := moofBox{children: make([]box, 0, len(children))}
	m.startPos = startPos
	for _, c := range children {
		err := m.addChild(c)
		if err != nil {
			return nil, err
		}
	}
	return &m, nil
}

// MoofBox -  Movie Fragment Box (moof)
//
// Contains all meta-data. To be able to stream a file, the moov box should be
// placed before the mdat box.
type moofBox struct {
	children []box
	startPos uint64
	Traf     *trafBox // The first traf child box
	Trafs    []*trafBox
}

// AddChild - add child box
func (m *moofBox) addChild(b box) error {
	switch b.getType() {
	case "traf":
		if m.Traf == nil {
			m.Traf = b.(*trafBox)
		}
		m.Trafs = append(m.Trafs, b.(*trafBox))
	}
	m.children = append(m.children, b)
	return nil
}

// Type - returns box type
func (m *moofBox) getType() string {
	return "moof"
}

// Size - returns calculated size
func (m *moofBox) size() uint64 {
	return containerSize(m.children)
}

// Encode - write moof after updating trun dataoffset
func (m *moofBox) encode(w io.Writer) error {
	for _, trun := range m.Traf.Truns {
		if trun.hasDataOffset() && trun.DataOffset == 0 {
			return fmt.Errorf("dataoffset in trun not set")
		}
	}
	err := encodeHeader(m, w)
	if err != nil {
		return err
	}
	for _, b := range m.children {
		err = b.encode(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetChildren - list of child boxes
func (m *moofBox) getChildren() []box {
	return m.children
}
