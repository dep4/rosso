package format

import (
   "io"
   "os"
   "strconv"
   "strings"
   "time"
   "unicode/utf8"
)

// mimesniff.spec.whatwg.org#binary-data-byte
func String(buf []byte) bool {
   for _, b := range buf {
      if b <= 0x08 {
         return false
      }
      if b == 0x0B {
         return false
      }
      if b >= 0x0E && b <= 0x1A {
         return false
      }
      if b >= 0x1C && b <= 0x1F {
         return false
      }
   }
   return utf8.Valid(buf)
}

type Number interface {
   float64 | int | int64 | ~uint64
}

func Label[T Number](value T, unit ...string) string {
   var (
      i int
      symbol string
      val = float64(value)
   )
   for i, symbol = range unit {
      if val < 1000 {
         break
      }
      val /= 1000
   }
   if i >= 1 {
      i = 3
   }
   return strconv.FormatFloat(val, 'f', i, 64) + symbol
}

func Label_Number[T Number](value T) string {
   return Label(value, "", " K", " M", " B", " T")
}

func Label_Rate[T Number](value T) string {
   return Label(value, " B/s", " kB/s", " MB/s", " GB/s", " TB/s")
}

func Label_Size[T Number](value T) string {
   return Label(value, " B", " kB", " MB", " GB", " TB")
}

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

func (p *Progress) Write(buf []byte) (int, error) {
   if p.time.IsZero() {
      p.time = time.Now()
      p.time_lap = time.Now()
   }
   since := time.Since(p.time_lap)
   if since >= time.Second {
      os.Stderr.WriteString(p.String())
      os.Stderr.WriteString("\n")
      p.time_lap = p.time_lap.Add(since)
   }
   write, err := p.w.Write(buf)
   p.bytes_written += write
   return write, err
}

func (p Progress) String() string {
   percent := func(value int, total int64) string {
      var ratio float64
      if total != 0 {
         ratio = 100 * float64(value) / float64(total)
      }
      return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
   }
   ratio := percent(p.bytes_written, p.bytes)
   rate := float64(p.bytes_written) / time.Since(p.time).Seconds()
   var buf strings.Builder
   buf.WriteString(ratio)
   buf.WriteByte('\t')
   buf.WriteString(Label_Size(p.bytes_written))
   buf.WriteByte('\t')
   buf.WriteString(Label_Rate(rate))
   return buf.String()
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