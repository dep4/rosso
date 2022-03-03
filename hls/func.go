package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
   "net/url"
   "strconv"
   "text/scanner"
   "unicode"
)

func (d Decrypter) Copy(dst io.Writer, src io.Reader) (int, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return 0, err
   }
   cipher.NewCBCDecrypter(d.Block, d.IV).CryptBlocks(buf, buf)
   if len(buf) >= 1 {
      pad := buf[len(buf)-1]
      if len(buf) >= int(pad) {
         buf = buf[:len(buf)-int(pad)]
      }
   }
   return dst.Write(buf)
}

func NewMaster(addr *url.URL, body io.Reader) (*Master, error) {
   var (
      buf scanner.Scanner
      err error
      mas Master
   )
   buf.Init(body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-MEDIA":
         var med Media
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "GROUP-ID":
               buf.Scan()
               buf.Scan()
               med.GroupID, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
            case "URI":
               buf.Scan()
               buf.Scan()
               med.URI, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               addr, err := addr.Parse(med.URI)
               if err != nil {
                  return nil, err
               }
               med.URI = addr.String()
            }
         }
         mas.Media = append(mas.Media, med)
      case "EXT-X-STREAM-INF":
         var str Stream
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               str.Resolution = buf.TokenText()
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               str.Bandwidth, err = strconv.ParseInt(buf.TokenText(), 10, 64)
            case "CODECS":
               buf.Scan()
               buf.Scan()
               str.Codecs, err = strconv.Unquote(buf.TokenText())
            case "AUDIO":
               buf.Scan()
               buf.Scan()
               str.Audio, err = strconv.Unquote(buf.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         scanLines(&buf)
         buf.Scan()
         addr, err := addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         str.URI = addr.String()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}

func NewSegment(addr *url.URL, body io.Reader) (*Segment, error) {
   var (
      buf scanner.Scanner
      err error
      seg Segment
   )
   buf.Init(body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-KEY":
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "METHOD":
               buf.Scan()
               buf.Scan()
               seg.Key.Method = buf.TokenText()
            case "URI":
               buf.Scan()
               buf.Scan()
               seg.Key.URI, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               addr, err := addr.Parse(seg.Key.URI)
               if err != nil {
                  return nil, err
               }
               seg.Key.URI = addr.String()
            }
         }
      case "EXTINF":
         var info Information
         buf.Scan()
         buf.Scan()
         info.Duration = buf.TokenText()
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         addr, err := addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         info.URI = addr.String()
         seg.Info = append(seg.Info, info)
      }
   }
   return &seg, nil
}

func scanLines(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r != '\n'
   }
   buf.Whitespace = 1 << '\n'
}

func scanWords(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r == '-' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
   }
   buf.Whitespace = 1 << ' '
}

func NewDecrypter(src io.Reader) (*Decrypter, error) {
   key, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Decrypter{block, key}, nil
}

func (s Stream) String() string {
   var buf []byte
   if s.Resolution != "" {
      buf = append(buf, "Resolution:"...)
      buf = append(buf, s.Resolution...)
      buf = append(buf, ' ')
   }
   buf = append(buf, "Bandwidth:"...)
   buf = strconv.AppendInt(buf, s.Bandwidth, 10)
   buf = append(buf, " Codecs:"...)
   buf = append(buf, s.Codecs...)
   if s.Audio != "" {
      buf = append(buf, " Audio:"...)
      buf = append(buf, s.Audio...)
   }
   if s.URI != "" {
      buf = append(buf, " URI:"...)
      buf = append(buf, s.URI...)
   }
   return string(buf)
}
