package os

import (
   "github.com/89z/rosso/strconv"
   "io"
   "os"
   "time"
)

var (
   Args = os.Args
   Link = os.Link
   Open = os.Open
   ReadFile = os.ReadFile
   Stat = os.Stat
   Stderr = os.Stderr
   Stdout = os.Stdout
   UserHomeDir = os.UserHomeDir
)

type Progress struct {
   bytes int64
   bytes_read int64
   bytes_written int
   chunks int
   chunks_read int64
   lap time.Time
   total time.Time
   w io.Writer
}

func Progress_Bytes(dst io.Writer, bytes int64) *Progress {
   return &Progress{w: dst, bytes: bytes}
}

func Progress_Chunks(dst io.Writer, chunks int) *Progress {
   return &Progress{w: dst, chunks: chunks}
}

func (p *Progress) Add_Chunk(bytes int64) {
   p.bytes_read += bytes
   p.chunks_read += 1
   p.bytes = int64(p.chunks) * p.bytes_read / p.chunks_read
}

func (p *Progress) Write(buf []byte) (int, error) {
   if p.total.IsZero() {
      p.total = time.Now()
      p.lap = time.Now()
   }
   lap := time.Since(p.lap)
   if lap >= time.Second {
      total := time.Since(p.total).Seconds()
      var b []byte
      b = strconv.NewRatio(p.bytes_written, p.bytes).AppendPercent(b)
      b = append(b, "   "...)
      b = strconv.AppendSize(b, p.bytes_written)
      b = append(b, "   "...)
      b = strconv.NewRatio(p.bytes_written, total).AppendRate(b)
      b = append(b, '\n')
      os.Stderr.Write(b)
      p.lap = p.lap.Add(lap)
   }
   write, err := p.w.Write(buf)
   p.bytes_written += write
   return write, err
}
