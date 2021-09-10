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

type Hashes []struct {
   MD5 string
   JA3 string
}

func NewHashes(r io.Reader) (Hashes, error) {
   var h Hashes
   if err := json.NewDecoder(r).Decode(&h); err != nil {
      return nil, err
   }
   return h, nil
}

func (h Hashes) JA3(md5 string) string {
   for _, hash := range h {
      if hash.MD5 == md5 {
         return hash.JA3
      }
   }
   return ""
}

type Users []struct {
   MD5 string
   Count int
   Agent string `json:"User-Agent"`
}

func NewUsers(r io.Reader) (Users, error) {
   var u Users
   if err := json.NewDecoder(r).Decode(&u); err != nil {
      return nil, err
   }
   return u, nil
}

func (u Users) Agents(md5 string) []string {
   var a []string
   for _, user := range u {
      if user.MD5 == md5 {
         a = append(a, user.Agent)
      }
   }
   return a
}

func (u Users) Sort() {
   sort.Slice(u, func(a, b int) bool {
      return u[b].Count < u[a].Count
   })
}
