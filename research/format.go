package format

import (
   "github.com/89z/format"
   "io"
   "os"
   "time"
)

type Progress struct {
   io.Writer
   length int
   lengthTotal int64
   time time.Time
   timeTotal time.Time
}

func NewProgress(src io.Writer, length int64) *Progress {
   var pro Progress
   pro.Writer = src
   pro.lengthTotal = length
   pro.time = time.Now()
   pro.timeTotal = time.Now()
   return &pro
}

func (p *Progress) Write(buf []byte) (int, error) {
   since := time.Since(p.time)
   if since >= time.Second/2 {
      p.progress()
      p.time = p.time.Add(since)
   }
   write, err := p.Writer.Write(buf)
   p.length += write
   return write, err
}

func (p Progress) progress() {
   rate := float64(p.length) / time.Since(p.timeTotal).Seconds()
   os.Stderr.WriteString(format.Percent(p.length, p.lengthTotal))
   os.Stderr.WriteString("\t")
   os.Stderr.WriteString(format.LabelSize(p.length))
   os.Stderr.WriteString("\t")
   os.Stderr.WriteString(format.LabelRate(rate))
   os.Stderr.WriteString("\n")
}
