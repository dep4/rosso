package format

import (
   "github.com/89z/format"
   "io"
   "net/http"
   "os"
   "strings"
   "time"
)

type Progress struct {
   io.Writer
   chunk struct {
      current int64
      total int64
   }
   length struct {
      previous int
      current int64
   }
   time struct {
      start time.Time
      lap time.Time
   }
}

func NewProgress(src io.Writer, chunks int) *Progress {
   var pro Progress
   pro.Writer = src
   pro.chunk.total = int64(chunks)
   pro.time.lap = time.Now()
   pro.time.start = time.Now()
   return &pro
}

func (p Progress) String() string {
   lengthTotal := p.length.current / p.chunk.current * p.chunk.total
   rate := float64(p.length.previous) / time.Since(p.time.start).Seconds()
   var buf strings.Builder
   buf.WriteString(format.Percent(p.length.previous, lengthTotal))
   buf.WriteByte('\t')
   buf.WriteString(format.LabelSize(p.length.previous))
   buf.WriteByte('\t')
   buf.WriteString(format.LabelRate(rate))
   return buf.String()
}

func (p *Progress) Copy(res *http.Response) (int64, error) {
   p.chunk.current += 1
   p.length.current += res.ContentLength
   return io.Copy(p, res.Body)
}

func (p *Progress) Write(buf []byte) (int, error) {
   since := time.Since(p.time.lap)
   if since >= time.Second/2 {
      os.Stderr.WriteString(p.String())
      os.Stderr.WriteString("\n")
      p.time.lap = p.time.lap.Add(since)
   }
   write, err := p.Writer.Write(buf)
   p.length.previous += write
   return write, err
}
