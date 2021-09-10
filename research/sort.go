package sort

import (
   "encoding/json"
   "io"
)

const (
   allHashes = "https://ja3er.com/getAllHashesJson"
   allUas = "https://ja3er.com/getAllUasJson"
)

type counts []struct {
   MD5 string
   UserAgent string `json:"User-Agent"`
}

func newCounts(r io.Reader) (counts, error) {
   var c counts
   if err := json.NewDecoder(r).Decode(&c); err != nil {
      return nil, err
   }
   return c, nil
}

func (c counts) agents(md5 string) []string {
   var a []string
   for _, count := range c {
      if count.MD5 == md5 {
         a = append(a, count.UserAgent)
      }
   }
   return a
}

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
