package format

import (
   "io"
   "net/http"
   "time"
)

type Progress struct {
   *http.Response
   content, partLength, part int64
   io.Writer
   time.Time
}

func NewProgress(src *http.Response, dst io.Writer) *Progress {
   var pro Progress
   pro.Response = src
   pro.Time = time.Now()
   pro.Writer = dst
   pro.partLength = 10_000_000
   return &pro
}

func (p *Progress) Read(buf []byte) (int, error) {
   if p.part == 0 {
      end := time.Since(p.Time).Milliseconds()
      if end >= 1 {
         PercentInt64(p.Writer, p.content, p.ContentLength)
         io.WriteString(p.Writer, "\t")
         Size.LabelInt64(p.Writer, p.content)
         io.WriteString(p.Writer, "\t")
         Rate.LabelInt64(p.Writer, 1000 * p.content / end)
         io.WriteString(p.Writer, "\n")
      }
   }
   // Callers should always process the n > 0 bytes returned before considering
   // the error err.
   read, err := p.Body.Read(buf)
   p.content += int64(read)
   p.part += int64(read)
   if p.part >= p.partLength {
      p.part = 0
   }
   return read, err
}
