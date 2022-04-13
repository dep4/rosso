package format

import (
   "io"
   "os"
   "strings"
   "time"
)

func (p *Progress) Write(buf []byte) (int, error) {
   since := time.Since(p.lapTime)
   if since >= time.Second/2 {
      os.Stderr.WriteString(p.String())
      os.Stderr.WriteString("\n")
      p.lapTime = p.lapTime.Add(since)
   }
   write, err := p.Writer.Write(buf)
   p.bytes += write
   return write, err
}

////////////////////////////////////////////////////////////////////////////////

func (p Progress) String() string {
   rate := float64(p.bytes) / time.Since(p.time).Seconds()
   var buf strings.Builder
   buf.WriteString(Percent(p.bytes, p.chunks*p.bytesRead/p.chunksRead))
   buf.WriteByte('\t')
   buf.WriteString(LabelSize(p.bytes))
   buf.WriteByte('\t')
   buf.WriteString(LabelRate(rate))
   return buf.String()
}

type Progress struct {
   io.Writer
   bytes int
   bytesRead int64
   chunks int64
   chunksRead int64
   time time.Time
   lapTime time.Time
}

func (p *Progress) AddChunk(length int64) {
   p.chunksRead += 1
   p.bytesRead += length
}

func NewProgress(dst io.Writer, chunks int) *Progress {
   var pro Progress
   pro.Writer = dst
   pro.chunks = int64(chunks)
   pro.lapTime = time.Now()
   pro.time = time.Now()
   return &pro
}
