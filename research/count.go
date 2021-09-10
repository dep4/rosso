package sort

import (
   "encoding/json"
   "fmt"
   "io"
   "sort"
)

func (c counts) groups() groups {
   var g groups
   m := make(map[string]int)
   for i, count := range c {
      fmt.Printf("%9v\r", len(c)-i)
      j, ok := m[count.MD5]
      if ok {
         g[j].agents = append(g[j].agents, count.UserAgent)
      } else {
         g = append(g, group{
            count.MD5, []string{count.UserAgent},
         })
         m[count.MD5] = len(g)-1
      }
   }
   return g
}

type group struct {
   md5 string
   agents []string
}

func (c counts) filter(lastSeen string) counts {
   var d counts
   for _, count := range c {
      if count.LastSeen > lastSeen {
         d = append(d, count)
      }
   }
   return d
}

type groups []group

func (g groups) filter(agents int) groups {
   var h groups
   for _, group := range g {
      if len(group.agents) > agents {
         h = append(h, group)
      }
   }
   return h
}

func (g groups) sort() {
   sort.Slice(g, func(a, b int) bool {
      ga, gb := g[a], g[b]
      return len(gb.agents) < len(ga.agents)
   })
}
