package hls

import (
   "container/heap"
   "fmt"
   "testing"
)

func TestSort(t *testing.T) {
   mas := &Master{Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 144},
      {Bandwidth: 1080},
      {Bandwidth: 720},
      {Bandwidth: 2160},
   }}
   heap.Init(mas)
   fmt.Printf("%+v\n", heap.Pop(mas))
}
