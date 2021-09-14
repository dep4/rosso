package ja3

import (
   "encoding/json"
   "io"
   "sort"
)

const (
   AllHashes = "https://ja3er.com/getAllHashesJson"
   AllUas = "https://ja3er.com/getAllUasJson"
)

type JA3er struct {
   Users []struct {
      MD5 string
      Count int
      Agent string `json:"User-Agent"`
   }
   Hashes []struct {
      MD5 string
      JA3 string
   }
}

func NewJA3er(ua, hash io.Reader) (*JA3er, error) {
   var j JA3er
   if err := json.NewDecoder(ua).Decode(&j.Users); err != nil {
      return nil, err
   }
   if err := json.NewDecoder(hash).Decode(&j.Hashes); err != nil {
      return nil, err
   }
   return &j, nil
}

func (j JA3er) JA3(md5 string) string {
   for _, hash := range j.Hashes {
      if hash.MD5 == md5 {
         return hash.JA3
      }
   }
   return ""
}

func (j JA3er) SortUsers() {
   sort.Slice(j.Users, func(a, b int) bool {
      return j.Users[b].Count < j.Users[a].Count
   })
}
