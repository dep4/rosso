package os

import (
   "github.com/89z/std/strconv"
   "io"
   "os"
   "path/filepath"
   "strings"
   "time"
)

func (p *Progress) Write(buf []byte) (int, error) {
   if p.time.IsZero() {
      p.time = time.Now()
      p.time_lap = time.Now()
   }
   since := time.Since(p.time_lap)
   if since >= time.Second {
      os.Stderr.WriteString(strconv.Percent(p.bytes_written, p.bytes))
      os.Stderr.WriteString("\t")
      os.Stderr.WriteString(strconv.Size(p.bytes_written))
      os.Stderr.WriteString("\t")
      os.Stderr.WriteString(strconv.Rate(p.bytes_written, p.time))
      os.Stderr.WriteString("\n")
      p.time_lap = p.time_lap.Add(since)
   }
   write, err := p.w.Write(buf)
   p.bytes_written += write
   return write, err
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

var (
   Open = os.Open
   ReadFile = os.ReadFile
   Stdout = os.Stdout
)

func Create(name string) (*os.File, error) {
   var err error
   name, err = clean(name)
   if err != nil {
      return nil, err
   }
   return os.Create(name)
}

func WriteFile(name string, data []byte) error {
   var err error
   name, err = clean(name)
   if err != nil {
      return err
   }
   return os.WriteFile(name, data, os.ModePerm)
}

func clean(name string) (string, error) {
   dir, file := filepath.Split(name)
   if dir != "" {
      err := os.MkdirAll(dir, os.ModePerm)
      if err != nil {
         return "", err
      }
   }
   mapping := func(r rune) rune {
      if strings.ContainsRune(`"*/:<>?\|`, r) {
         return -1
      }
      return r
   }
   file = strings.Map(mapping, file)
   name = filepath.Join(dir, file)
   os.Stderr.WriteString("OpenFile " + name + "\n")
   return name, nil
}
