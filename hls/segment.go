package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
   "net/url"
   "strconv"
   "text/scanner"
)

type Decrypter struct {
   cipher.Block
   IV []byte
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

// We do not care about the ciphertext, so this works in place.
func (d Decrypter) Decrypt(src io.Reader) ([]byte, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   cipher.NewCBCDecrypter(d.Block, d.IV).CryptBlocks(buf, buf)
   if len(buf) >= 1 {
      pad := buf[len(buf)-1]
      if len(buf) >= int(pad) {
         buf = buf[:len(buf)-int(pad)]
      }
   }
   return buf, nil
}

type Information struct {
   Duration string
   URI string
}

type Segment struct {
   Key struct {
      Method string
      URI string
   }
   Info []Information
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
      case "EXTINF":
         var info Information
         buf.Scan()
         buf.Scan()
         info.Duration = buf.TokenText()
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         addr, err = addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         info.URI = addr.String()
         seg.Info = append(seg.Info, info)
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
            }
         }
      }
   }
   return &seg, nil
}
