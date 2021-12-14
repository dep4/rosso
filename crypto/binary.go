package crypto

import (
   "encoding/binary"
   "github.com/89z/parse"
   "github.com/refraction-networking/utls"
   "io"
)

func ParseHandshake(data []byte) (*ClientHello, error) {
   // Content Type, Version
   if dLen := len(data); dLen <= 2 {
      return nil, parse.InvalidSlice{2, dLen}
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

func extensionType(ext tls.TLSExtension) (uint16, error) {
   data, err := io.ReadAll(ext)
   if err != nil || len(data) <= 1 {
      return 0, err
   }
   return binary.BigEndian.Uint16(data), nil
}
