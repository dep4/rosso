package format

import (
   "io"
   "os"
   "strings"
   "time"
)

type Progress struct {
   io.Writer
   bytes int64
   bytesRead int64
   bytesWritten int
   chunks int
   chunksRead int64
   time time.Time
   timeLap time.Time
}

func ProgressBytes(dst io.Writer, bytes int64) *Progress {
   return &Progress{Writer: dst, bytes: bytes}
}

func ProgressChunks(dst io.Writer, chunks int) *Progress {
   return &Progress{Writer: dst, chunks: chunks}
}

func (p *Progress) AddChunk(bytes int64) {
   p.bytesRead += bytes
   p.chunksRead += 1
   p.bytes = int64(p.chunks) * p.bytesRead / p.chunksRead
}

func (p Progress) String() string {
   rate := float64(p.bytesWritten) / time.Since(p.time).Seconds()
   var buf strings.Builder
   buf.WriteString(Percent(p.bytesWritten, p.bytes))
   buf.WriteByte('\t')
   buf.WriteString(LabelSize(p.bytesWritten))
   buf.WriteByte('\t')
   buf.WriteString(LabelRate(rate))
   return buf.String()
}

func (p *Progress) Write(buf []byte) (int, error) {
   if p.time.IsZero() {
      p.time = time.Now()
      p.timeLap = time.Now()
   }
   since := time.Since(p.timeLap)
   if since >= time.Second/2 {
      os.Stderr.WriteString(p.String())
      os.Stderr.WriteString("\n")
      p.timeLap = p.timeLap.Add(since)
   }
   write, err := p.Writer.Write(buf)
   p.bytesWritten += write
   return write, err
}
