package m3u

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "io"
   "net/http"
   "strconv"
)

func unpad(buf []byte) []byte {
   total := len(buf)
   if total > 0 {
      value := int(buf[total-1])
      if value < total {
         return buf[:total-value]
      }
   }
   return nil
}

func (f Format) blockMode() (cipher.BlockMode, error) {
   res, err := http.Get(f["URI"])
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   key, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   // The CBS key is 16 bytes, which means BlockSize will be 16.
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return cipher.NewCBCDecrypter(block, key), nil
}

type Format map[string]string

func Unmarshal(buf []byte) []Format {
   lines := bytes.FieldsFunc(buf, func(r rune) bool {
      return r == '\n'
   })
   var pass1 []Format
   for _, line := range lines {
      if line[0] == '#' {
         form := make(Format)
         pairs := reader{line}
         pairs.readBytes(':', '"')
         for {
            if pairs.buf == nil {
               break
            }
            var pair reader
            pair.buf = pairs.readBytes(',', '"')
            key := pair.readBytes('=', '"')
            if pair.buf != nil {
               val := string(pair.buf)
               unq, err := strconv.Unquote(val)
               if err == nil {
                  val = unq
               }
               form[string(key)] = val
            }
         }
         pass1 = append(pass1, form)
      } else {
         ind := merge(pass1)
         if ind >= 0 {
            pass1[ind]["URI"] = string(line)
         } else {
            form := make(Format)
            form["URI"] = string(line)
            pass1 = append(pass1, form)
         }
      }
   }
   var pass2 []Format
   uris := make(map[string]bool)
   for _, form := range pass1 {
      uri, ok := form["URI"]
      if ok && !uris[uri] {
         form["URI"] = form["URI"]
         pass2 = append(pass2, form)
         uris[uri] = true
      }
   }
   return pass2
}

func merge(forms []Format) int {
   if len(forms) >= 1 {
      form := forms[len(forms)-1]
      if len(form) >= 1 {
         return len(forms)-1
      }
   }
   return -1
}
