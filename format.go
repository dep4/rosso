package format

import (
   "bytes"
   "fmt"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
   "time"
)

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
)

// godocs.io/github.com/google/pprof/internal/measurement#Percentage
func Percent(value, total float64) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

// godocs.io/github.com/google/pprof/internal/measurement#Percentage
func PercentInt(value, total int) string {
   val, tot := float64(value), float64(total)
   return Percent(val, tot)
}

// godocs.io/github.com/google/pprof/internal/measurement#Percentage
func PercentInt64(value, total int64) string {
   val, tot := float64(value), float64(total)
   return Percent(val, tot)
}

func Trim(s string) string {
   if len(s) <= 99 {
      return s
   }
   return s[:48] + "..." + s[len(s)-48:]
}

type LogLevel int

func (l LogLevel) Dump(req *http.Request) error {
   switch l {
   case 0:
      loc := Trim(req.URL.String())
      fmt.Println(req.Method, loc)
   case 1:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return err
      }
      os.Stdout.Write(buf)
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         os.Stdout.WriteString("\n")
      }
   case 2:
      buf, err := httputil.DumpRequestOut(req, true)
      if err != nil {
         return err
      }
      os.Stdout.Write(buf)
   }
   return nil
}

type Progress struct {
   *http.Response
   Content, PartLength, part int64
   time.Time
}

func Content(length int64) Progress {
   pro := Progress{PartLength: 10_000_000}
   pro.ContentLength = length
   pro.Time = time.Now()
   return pro
}

func Response(res *http.Response) *Progress {
   pro := Progress{Response: res, PartLength: 10_000_000}
   pro.Time = time.Now()
   return &pro
}

func (p Progress) Print() {
   end := time.Since(p.Time).Milliseconds()
   if end > 0 {
      meter := PercentInt64(p.Content, p.ContentLength)
      meter += "\t" + Size.LabelInt(p.Content)
      meter += "\t" + Rate.LabelInt(1000 * p.Content / end)
      fmt.Println(meter)
   }
}

func (p Progress) Range() string {
   buf := []byte("bytes=")
   buf = strconv.AppendInt(buf, p.Content, 10)
   buf = append(buf, '-')
   buf = strconv.AppendInt(buf, p.Content+p.PartLength-1, 10)
   return string(buf)
}

func (p *Progress) Read(buf []byte) (int, error) {
   if p.part == 0 {
      p.Print()
   }
   read, err := p.Body.Read(buf)
   if err != nil {
      return 0, err
   }
   p.Content += int64(read)
   p.part += int64(read)
   if p.part >= p.PartLength {
      p.part = 0
   }
   return read, nil
}

type Symbols []string

// godocs.io/github.com/google/pprof/internal/measurement#Label
func (s Symbols) Label(f float64) string {
   var (
      i int
      symbol string
   )
   for i, symbol = range s {
      if f < 1000 {
         break
      }
      f /= 1000
   }
   if i == 0 {
      return strconv.FormatFloat(f, 'f', 0, 64) + symbol
   }
   return strconv.FormatFloat(f, 'f', 3, 64) + symbol
}

// godocs.io/github.com/google/pprof/internal/measurement#Label
func (s Symbols) LabelInt(i int64) string {
   f := float64(i)
   return s.Label(f)
}

// godocs.io/github.com/google/pprof/internal/measurement#Label
func (s Symbols) LabelUint(i uint64) string {
   f := float64(i)
   return s.Label(f)
}
