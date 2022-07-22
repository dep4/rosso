package os

import (
   "github.com/89z/rosso/strconv"
   "io"
   "os"
   "path/filepath"
   "strings"
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

func Create(name string) (*os.File, error) {
   err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
   if err != nil {
      return nil, err
   }
   os.Stderr.WriteString("Create " + name + "\n")
   return os.Create(name)
}

func Rename(old_path, new_path string) error {
   err := os.MkdirAll(filepath.Dir(new_path), os.ModePerm)
   if err != nil {
      return err
   }
   os.Stderr.WriteString("Rename " + new_path + "\n")
   return os.Rename(old_path, new_path)
}

func WriteFile(name string, data []byte) error {
   err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
   if err != nil {
      return err
   }
   os.Stderr.WriteString("WriteFile " + name + "\n")
   return os.WriteFile(name, data, os.ModePerm)
}

type Cleaner struct {
   dir string
   name string
}

func Clean(dir, file string) Cleaner {
   mapping := func(r rune) rune {
      if strings.ContainsRune(`"*/:<>?\|`, r) {
         return -1
      }
      return r
   }
   var c Cleaner
   c.dir = dir
   c.name = strings.Map(mapping, file)
   c.name = filepath.Join(c.dir, c.name)
   return c
}

func (c Cleaner) Create() (*os.File, error) {
   err := os.MkdirAll(c.dir, os.ModePerm)
   if err != nil {
      return nil, err
   }
   os.Stderr.WriteString("Create " + c.name + "\n")
   return os.Create(c.name)
}

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
