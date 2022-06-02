package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

func (h boxHeader) readBoxBody(r io.Reader) ([]byte, error) {
	bodyLen := h.size - uint64(h.length)
	if bodyLen == 0 {
		return nil, nil
	}
	body := make([]byte, bodyLen)
	_, err := io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// BoxHeader - 8 or 16 bytes depending on size
type boxHeader struct {
	name   string
	size   uint64
	length int
}

func decodeHeader(r io.Reader) (*boxHeader, error) {
	buf := make([]byte, boxHeaderSize)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	var head boxHeader
	head.size = uint64(binary.BigEndian.Uint32(buf))
	head.length = boxHeaderSize
	if head.size == 1 {
		buf := make([]byte, largeSizeLen)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}
		head.size = binary.BigEndian.Uint64(buf)
		head.length += largeSizeLen
	} else if head.size == 0 {
		return nil, fmt.Errorf("size 0, meaning to end of file, not supported")
	}
	head.name = string(buf[4:8])
	return &head, nil
}

func decodeBox(startPos uint64, r io.Reader) (box, error) {
	h, err := decodeHeader(r)
	if err != nil {
		return nil, err
	}
	var b box
	switch h.name {
	case "mdat":
		b, err = decodeMdat(h, startPos, r)
	case "moof":
		b, err = decodeMoof(h, startPos, r)
	default:
		b, err = decodeUnknown(h, startPos, r)
	}
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", h.name, err)
	}
	return b, nil
}

// ContainerBox is interface for ContainerBoxes
type containerBox interface {
	encode(w io.Writer) error
	getChildren() []box
	size() uint64
	getType() string
}

func containerSize(children []box) uint64 {
	var contentSize uint64
	for _, child := range children {
		contentSize += child.size()
	}
	return boxHeaderSize + contentSize
}

// DecodeContainerChildren decodes a container box
func decodeContainerChildrenSR(hdr *boxHeader, startPos, endPos uint64, sr sliceReader) ([]box, error) {
	var children []box
	pos := startPos
	initPos := sr.GetPos()
	for {
		if pos > endPos {
			return nil, fmt.Errorf("non matching children box sizes")
		}
		if pos == endPos {
			break
		}
		child, err := decodeBoxSR(pos, sr)
		if err != nil {
			return children, err
		}
		children = append(children, child)
		pos += child.size()
		relPosFromSize := sr.GetPos() - initPos
		if int(pos-startPos) != relPosFromSize {
			return nil, fmt.Errorf("child %s size mismatch in %s: %d - %d", child.getType(), hdr.name, pos-startPos, relPosFromSize)
		}
	}
	return children, nil
}

func encodeContainer(c containerBox, w io.Writer) error {
	err := encodeHeader(c, w)
	if err != nil {
		return err
	}
	for _, child := range c.getChildren() {
		err := child.encode(w)
		if err != nil {
			return err
		}
	}
	return nil
}

type box interface {
	// Type of box, normally 4 asccii characters, but is uint32 according to spec
	getType() string
	// Size of box including header and all children if any
	size() uint64
	// Encode box to writer
	encode(w io.Writer) error
}

const (
	// boxHeaderSize - standard size + name header
	boxHeaderSize = 8
	largeSizeLen  = 8          // Length of largesize exension
	flagsMask     = 0x00ffffff // Flags for masks from full header
)

func (b boxHeader) payloadLen() int {
	return int(b.size) - b.length
}

func encodeHeader(b box, w io.Writer) error {
	buf := make([]byte, boxHeaderSize)
	boxType := b.getType()
	strtobuf(buf[4:], boxType, 4)
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func encodeHeaderSW(b box, sw sliceWriter) {
	boxSize := b.size()
	sw.WriteUint32(uint32(boxSize))
	boxType := b.getType()
	sw.WriteString(boxType, false)
}

func strtobuf(out []byte, in string, l int) {
	if l < len(in) {
		copy(out, in)
	} else {
		copy(out, in[0:l])
	}
}
