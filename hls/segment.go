package hls

import (
   "encoding/hex"
   "strconv"
   "strings"
   "text/scanner"
)

func (s Scanner) Segment() (*Segment, error) {
   var (
      key bool
      seg Segment
   )
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      switch {
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         key = true
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               seg.Raw_IV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               var err error
               seg.Key, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         if key {
            seg.Protected = append(seg.Protected, line)
         } else {
            seg.Clear = append(seg.Clear, line)
         }
      case line == "#EXT-X-DISCONTINUITY":
         key = false
      }
   }
   return &seg, nil
}

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(s.Raw_IV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

type Segment struct {
   Clear []string
   Key string
   Protected []string
   Raw_IV string
}
