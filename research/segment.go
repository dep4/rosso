package hls

import (
   "encoding/hex"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
)

func (i Information) IV() ([]byte, error) {
   up := strings.ToUpper(i.RawIV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
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
      info Information
      seg Segment
   )
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      s.Init(strings.NewReader(line))
      switch {
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               info.RawIV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               var err error
               seg.RawKey, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         info.RawURI = line
         seg.Info = append(seg.Info, info)
         info = Information{}
      }
   }
   return &seg, nil
}

type Information struct {
   RawIV string
   RawURI string
}
