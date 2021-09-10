package ja3

import (
   "encoding/json"
   "io"
)

const AllUas = "https://ja3er.com/getAllUasJson"

type Group struct {
   MD5 string
   Agents []string
}

type Users []struct {
   MD5 string
   Agent string `json:"User-Agent"`
}

func NewUsers(r io.Reader) (Users, error) {
   var u Users
   if err := json.NewDecoder(r).Decode(&u); err != nil {
      return nil, err
   }
   return u, nil
}

func (u Users) Groups() []Group {
   var g []Group
   m := make(map[string]int)
   for _, user := range u {
      i, ok := m[user.MD5]
      if ok {
         g[i].Agents = append(g[i].Agents, user.Agent)
      } else {
         g = append(g, Group{
            user.MD5, []string{user.Agent},
         })
      }
   }
   return g
}
