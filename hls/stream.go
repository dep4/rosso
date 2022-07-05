package hls

type Stream struct {
   Audio string
   Bandwidth int
   Codecs string
   Resolution string
   URI string
}

type Stream_Filter interface {
   Audio(string) bool
   Bandwidth(int) int
   Codecs(string) bool
}

func (s Streams) Streams(f Stream_Filter) Streams {
   var prev Streams
   for _, curr := range s {
      if !f.Audio(curr.Audio) {
         continue
      }
      if !f.Codecs(curr.Codecs) {
         continue
      }
      if f.Bandwidth(curr.Bandwidth) < 0 {
         continue
      }
      prev = append(prev, curr)
   }
   return prev
}

func (s Streams) Bandwidth(f Stream_Filter) *Stream {
   var prev *Stream
   for i, curr := range s {
      if prev != nil {
         if f.Bandwidth(curr.Bandwidth) >= f.Bandwidth(prev.Bandwidth) {
            continue
         }
      }
      prev = &s[i]
   }
   return prev
}
