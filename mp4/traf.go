package mp4

import (
	"fmt"
	"io"
)

// DONE
func decodeFile(r io.Reader) (*file, error) {
	var (
		boxStartPos uint64
		f           file
	)
	for {
		box, err := decodeBox(boxStartPos, r)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		moof, ok := box.(*moofBox)
		if ok {
			for _, traf := range moof.Trafs {
				if traf.Senc == nil {
					return nil, fmt.Errorf("no senc box")
				}
				err := traf.Senc.parseReadBox()
				if err != nil {
					return nil, err
				}
			}
		}
		f.addChild(box, boxStartPos)
		boxStartPos += box.size()
	}
	return &f, nil
}

// AddChild - add child with start position
func (f *file) addChild(boxVal box, boxStartPos uint64) {
   switch boxVal.getType() {
   case "moof":
      f.IsFragmented = true
      moof := boxVal.(*moofBox)
      moof.startPos = boxStartPos
      var currentSegment *mediaSegment
      if len(f.Segments) == 0 {
         // No styp present, so one fragment per segment
         currentSegment = new(mediaSegment)
         f.addMediaSegment(currentSegment)
      } else {
         currentSegment = f.lastSegment()
      }
      frag := new(fragment)
      frag.addChild(moof)
      currentSegment.Fragments = append(currentSegment.Fragments, frag)
   case "mdat":
      mdat := boxVal.(*mdatBox)
      frags := f.lastSegment().Fragments
      currentFragment := frags[len(frags)-1]
      currentFragment.addChild(mdat)
   }
}

type file struct {
	Segments     []*mediaSegment // Media segments
	IsFragmented bool
}

// EncOptimize - encoder optimization mode
type encOptimize uint32

// AddMediaSegment - add a mediasegment to file f
func (f *file) addMediaSegment(m *mediaSegment) {
	f.Segments = append(f.Segments, m)
}

// LastSegment - Currently last segment
func (f *file) lastSegment() *mediaSegment {
	return f.Segments[len(f.Segments)-1]
}

type trafBox struct {
	children []box
	Senc     *sencBox
	Tfdt     *tfdtBox
	Tfhd     *tfhdBox
	Trun     *trunBox // The first TrunBox
	Truns    []*trunBox
}

// DecodeTrafSR - box-specific decode
func decodeTrafSR(hdr *boxHeader, startPos uint64, sr sliceReader) (box, error) {
	children, err := decodeContainerChildrenSR(hdr, startPos+8, startPos+hdr.size, sr)
	if err != nil {
		return nil, err
	}
	t := &trafBox{children: make([]box, 0, len(children))}
	for _, child := range children {
		err := t.addChild(child)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// AddChild - add child box
func (t *trafBox) addChild(b box) error {
	switch b.getType() {
	case "tfhd":
		t.Tfhd = b.(*tfhdBox)
	case "tfdt":
		t.Tfdt = b.(*tfdtBox)
	case "senc":
		t.Senc = b.(*sencBox)
	case "trun": // NEED THIS
		if t.Trun == nil {
			t.Trun = b.(*trunBox)
		}
		t.Truns = append(t.Truns, b.(*trunBox))
	}
	t.children = append(t.children, b)
	return nil
}

// Type - return box type
func (t *trafBox) getType() string {
	return "traf"
}

// Size - return calculated size
func (t *trafBox) size() uint64 {
	return containerSize(t.children)
}

// GetChildren - list of child boxes
func (t *trafBox) getChildren() []box {
	return t.children
}

// Encode - write box to w
func (t *trafBox) encode(w io.Writer) error {
	return encodeContainer(t, w)
}

// RemoveEncryptionBoxes - remove encryption boxes and return number of bytes
// removed
func (t *trafBox) removeEncryptionBoxes() uint64 {
	remainingChildren := make([]box, 0, len(t.children))
	var nrBytesRemoved uint64 = 0
	for _, ch := range t.children {
		switch ch.getType() {
		case "saiz":
			nrBytesRemoved += ch.size()
		case "saio":
			nrBytesRemoved += ch.size()
		case "senc":
			nrBytesRemoved += ch.size()
			t.Senc = nil
		case "uuid":
			nrBytesRemoved += ch.size()
		default:
			remainingChildren = append(remainingChildren, ch)
		}
	}
	t.children = remainingChildren
	return nrBytesRemoved
}
