package format

import (
   "fmt"
   "net/http"
   "time"
)

type Progress struct {
   *http.Response
   content int64
   part, partLength time.Time
}

func NewProgress(src *http.Response) *Progress {
   var pro Progress
   pro.Response = src
   pro.part = time.Now()
   pro.partLength = time.Now()
   return &pro
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
