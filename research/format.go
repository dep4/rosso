package format

import (
   "fmt"
   "net/http"
   "time"
)

type Progress struct {
   io.Writer
   time struct {
      value time.Time
      total time.Time
   }
   length struct {
      value int64
      total int64
   }
}

func NewProgress(src io.Writer, length int64) *Progress {
   var p Progress
   p.Writer = src
   p.length.total = length
   p.time.total = time.Now()
   p.time.value = time.Now()
   return &p
}

func (p *Progress) Read(buf []byte) (int, error) {
   since := time.Since(p.part)
   if since >= time.Second/2 {
      fmt.Println(p.content, p.ContentLength, time.Since(p.partLength))
      p.part = p.part.Add(since)
   }
   read, err := p.Body.Read(buf)
   p.content += int64(read)
   return read, err
}
