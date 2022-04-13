package format

import (
   "io"
   "os"
   "strings"
   "time"
)

type Progress struct {
   chunks int64
   io.Writer
   lapTime time.Time
   length int
   readChunks int64
   readLength int64
   time time.Time
}

func (p *Progress) AddChunk(length int64) {
   p.readChunks += 1
   p.readLength += length
}

func (p Progress) String() string {
   rate := float64(p.length) / time.Since(p.time).Seconds()
   var buf strings.Builder
   buf.WriteString(Percent(p.length, p.chunks*p.readLength/p.readChunks))
   buf.WriteByte('\t')
   buf.WriteString(LabelSize(p.length))
   buf.WriteByte('\t')
   buf.WriteString(LabelRate(rate))
   return buf.String()
}

func (p *Progress) Write(buf []byte) (int, error) {
   since := time.Since(p.lapTime)
   if since >= time.Second/2 {
      os.Stderr.WriteString(p.String())
      os.Stderr.WriteString("\n")
      p.lapTime = p.lapTime.Add(since)
   }
   write, err := p.Writer.Write(buf)
   p.length += write
   return write, err
}

func NewProgress(dst io.Writer, chunks int) *Progress {
   var pro Progress
   pro.Writer = dst
   pro.chunks = int64(chunks)
   pro.lapTime = time.Now()
   pro.time = time.Now()
   return &pro
}
