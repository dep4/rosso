package sort

import (
   "encoding/json"
   "io"
)

const allUas = "https://ja3er.com/getAllUasJson"

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

type group struct {
   md5 string
   agents []string
}

type groups []group

func (g groups) index(md5 string) int {
   for k, v := range g {
      if v.md5 == md5 {
         return k
      }
   }
   return -1
}
