package hls

import (
   "bytes"
   "encoding/hex"
   "net/url"
   "strings"
   "text/scanner"
)

func isInf(s []byte) bool {
   prefix := []byte("#EXTINF:")
   return bytes.HasPrefix(s, prefix)
}

func isKey(s []byte) bool {
   prefix := []byte("#EXT-X-KEY:")
   return bytes.HasPrefix(s, prefix)
}

func scanHex(s string) ([]byte, error) {
   up := strings.ToUpper(s)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

type Information struct {
   IV []byte
   URI *url.URL
}

func (s *Scanner) Segment(base *url.URL) (*Segment, error) {
   var (
      err error
      info Information
      seg Segment
   )
   for s.bufio.Scan() {
      slice := s.bufio.Bytes()
      s.Init(bytes.NewReader(slice))
      switch {
      case isKey(slice):
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               info.IV, err = scanHex(s.TokenText())
            case "URI":
               s.Scan()
               s.Scan()
               seg.Key, err = scanURL(base, s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
      case isInf(slice):
         s.bufio.Scan()
         info.URI, err = base.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
         seg.Info = append(seg.Info, info)
         info = Information{}
      }
   }
   return &seg, nil
}

type Segment struct {
   Key *url.URL
   Info []Information
}
