package sort

import (
   "encoding/json"
   "io"
)

const allHashes = "https://ja3er.com/getAllHashesJson"

type hashes []struct {
   JA3 string
   MD5 string
}

func newHashes(r io.Reader) (hashes, error) {
   var h hashes
   if err := json.NewDecoder(r).Decode(&h); err != nil {
      return nil, err
   }
   return h, nil
}

func (h hashes) ja3(md5 string) string {
   for _, hash := range h {
      if hash.MD5 == md5 {
         return hash.JA3
      }
   }
   return ""
}
