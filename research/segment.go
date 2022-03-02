package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "net/http"
   "strconv"
   "text/scanner"
)

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

func NewSegment(res *http.Response) (*Segment, error) {
   var (
      buf scanner.Scanner
      err error
      seg Segment
   )
   buf.Init(res.Body)
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
         addr, err := res.Request.URL.Parse(buf.TokenText())
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

type Block struct {
   cipher.Block
   IV []byte
}

func NewCipher(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{block, key}, nil
}

// We do not care about the ciphertext, so this works in place.
func (b Block) Decrypt(src []byte) []byte {
   total := len(src)
   if total >= b.BlockSize() {
      cipher.NewCBCDecrypter(b.Block, b.IV).CryptBlocks(src, src)
      value := int(src[total-1])
      if value < total {
         return src[:total-value]
      }
   }
   return nil
}
