package hls

import (
   "bytes"
   "encoding/hex"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
)

func isURI(s []byte) bool {
   prefix := []byte{'#'}
   return len(s) >= 1 && !bytes.HasPrefix(s, prefix)
}

func isKey(s []byte) bool {
   prefix := []byte("#EXT-X-KEY:")
   return bytes.HasPrefix(s, prefix)
}

func (i Information) IV() ([]byte, error) {
   up := strings.ToUpper(i.iv)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

type Information struct {
   iv string
   RawURI string
}

func (i Information) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(i.RawURI)
}

type Segment struct {
   Info []Information
   RawKey string
}

func (s Segment) Key(base *url.URL) (*url.URL, error) {
   return base.Parse(s.RawKey)
}

func (s Scanner) Segment() (*Segment, error) {
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
               info.iv = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               seg.RawKey, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case isURI(slice):
         info.RawURI = s.bufio.Text()
         seg.Info = append(seg.Info, info)
         info = Information{}
      }
   }
   return &seg, nil
}

