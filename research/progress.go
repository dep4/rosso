package format

import (
   "github.com/89z/format"
   "io"
   "os"
   "strings"
   "time"
)

var (
   LabelRate = format.LabelRate[float64]
   LabelSize = format.LabelSize[int]
   Percent = format.Percent[int, int64]
)

func newProgress(dst io.Writer) *progress {
   var p progress
   p.Writer = dst
   p.time.lap = time.Now()
   p.time.start = time.Now()
   return &p
}

func (p *progress) Write(buf []byte) (int, error) {
   since := time.Since(p.time.lap)
   if since >= time.Second/2 {
      os.Stderr.WriteString(p.String())
      os.Stderr.WriteString("\n")
      p.time.lap = p.time.lap.Add(since)
   }
   write, err := p.Writer.Write(buf)
   p.bytes.written += write
   return write, err
}

func (p *progress) setBytes(length int64) {
   p.bytes.total = length
}

func (p *progress) setChunks(length int64) {
   p.chunk.total = length
}

type progress struct {
   io.Writer
   time struct {
      start time.Time
      lap time.Time
   }
   chunk struct {
      read int64
      total int64
   }
   bytes struct {
      written int
      read int64
      total int64
   }
}

func (p *progress) addChunk(length int64) {
   p.bytes.read += length
   p.chunk.read += 1
   p.bytes.total = p.chunk.total * p.bytes.read / p.chunk.read
}

func (p progress) String() string {
   rate := float64(p.bytes.written) / time.Since(p.time.start).Seconds()
   var buf strings.Builder
   buf.WriteString(Percent(p.bytes.written, p.bytes.total))
   buf.WriteByte('\t')
   buf.WriteString(LabelSize(p.bytes.written))
   buf.WriteByte('\t')
   buf.WriteString(LabelRate(rate))
   return buf.String()
}
