package hls

import (
   "net/url"
)

type Media struct {
   GroupID string
   URI *url.URL
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle missing resolution
   Audio string // link to Media
   URI *url.URL
}

type Master struct {
   Stream []Stream
   Media []Media
}

func (m Master) Len() int {
   return len(m.Stream)
}

func (m *Master) Swap(i, j int) {
   m.Stream[i], m.Stream[j] = m.Stream[j], m.Stream[i]
}

func (m *Master) Pop() interface{} {
   high := len(m.Stream) - 1
   pop := m.Stream[high]
   m.Stream = m.Stream[:high]
   return pop
}

func (m *Master) Push(x interface{}) {
   m.Stream = append(m.Stream, x.(Stream))
}

func (m Master) Less(i, j int) bool {
   return m.Stream[i].Bandwidth < m.Stream[j].Bandwidth
}
