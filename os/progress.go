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

func (self *Progress) Add_Chunk(bytes int64) {
   self.bytes_read += bytes
   self.chunks_read += 1
   self.bytes = int64(self.chunks) * self.bytes_read / self.chunks_read
}

func (self *Progress) Write(buf []byte) (int, error) {
   if self.time.IsZero() {
      self.time = time.Now()
      self.time_lap = time.Now()
   }
   since := time.Since(self.time_lap)
   if since >= time.Second {
      var s string
      s += strconv.Percent(self.bytes_written, self.bytes)
      s += "\t"
      s += strconv.Size(self.bytes_written)
      s += "\t"
      s += strconv.Rate(self.bytes_written, time.Since(self.time).Seconds())
      s += "\n"
      os.Stderr.WriteString(s)
      self.time_lap = self.time_lap.Add(since)
   }
   write, err := self.w.Write(buf)
   self.bytes_written += write
   return write, err
}
