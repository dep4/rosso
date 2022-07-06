package os

import (
   "github.com/89z/rosso/strconv"
   "io"
   "os"
   "time"
)

type Progress struct {
   bytes int64
   bytes_read int64
   bytes_written int
   chunks int
   chunks_read int64
   time time.Time
   time_lap time.Time
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
   if p.time.IsZero() {
      p.time = time.Now()
      p.time_lap = time.Now()
   }
   since := time.Since(p.time_lap)
   if since >= time.Second {
      var s string
      s += strconv.Percent(p.bytes_written, p.bytes)
      s += "\t"
      s += strconv.Size(p.bytes_written)
      s += "\t"
      s += strconv.Rate(p.bytes_written, time.Since(p.time).Seconds())
      s += "\n"
      os.Stderr.WriteString(s)
      p.time_lap = p.time_lap.Add(since)
   }
   write, err := p.w.Write(buf)
   p.bytes_written += write
   return write, err
}
