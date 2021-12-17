package crypto

import (
   "encoding/binary"
   "github.com/refraction-networking/utls"
   "io"
   "strconv"
)

const AndroidJA3 =
   "769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-" +
   "51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

func extensionType(ext tls.TLSExtension) (uint16, error) {
   data, err := io.ReadAll(ext)
   if err != nil || len(data) <= 1 {
      return 0, err
   }
   return binary.BigEndian.Uint16(data), nil
}

func (b *Buffer) ReadUint16LengthPrefixed() ([]byte, []byte, bool) {
   low := 2
   if len(b.buf) < low {
      return nil, nil, false
   }
   high := low + int(binary.BigEndian.Uint16(b.buf))
   if len(b.buf) < high {
      return nil, nil, false
   }
   pre, buf := b.buf[:low], b.buf[low:high]
   b.buf = b.buf[high:]
   return pre, buf, true
}

// github.com/golang/go/issues/49227
func (b *Buffer) ReadUint32LengthPrefixed() ([]byte, []byte, bool) {
   low := 4
   if len(b.buf) < low {
      return nil, nil, false
   }
   high := low + int(binary.BigEndian.Uint32(b.buf))
   if len(b.buf) < high {
      return nil, nil, false
   }
   pre, buf := b.buf[:low], b.buf[low:high]
   b.buf = b.buf[high:]
   return pre, buf, true
}

func ParseHandshake(data []byte) (*ClientHello, error) {
   // Content Type, Version
   if dLen := len(data); dLen <= 2 {
      return nil, invalidSlice{2, dLen}
   }
   version := binary.BigEndian.Uint16(data[1:])
   // unsupported extension 0x16
   fin := tls.Fingerprinter{AllowBluntMimicry: true}
   spec, err := fin.FingerprintClientHello(data)
   if err != nil {
      return nil, err
   }
   return &ClientHello{spec, version}, nil
}

type invalidSlice struct {
   index, length int
}

func (i invalidSlice) Error() string {
   index, length := int64(i.index), int64(i.length)
   var buf []byte
   buf = append(buf, "index out of range ["...)
   buf = strconv.AppendInt(buf, index, 10)
   buf = append(buf, "] with length "...)
   buf = strconv.AppendInt(buf, length, 10)
   return string(buf)
}
