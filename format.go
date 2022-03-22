package format

import (
   "bytes"
   "encoding/json"
   "net/http"
   "net/http/httputil"
   "os"
   "path/filepath"
   "strconv"
   "time"
)

func Decode[T any](elem ...string) (*T, error) {
   name := filepath.Join(elem...)
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   value := new(T)
   if err := json.NewDecoder(file).Decode(value); err != nil {
      return nil, err
   }
   return value, nil
}

func Encode[T any](value T, elem ...string) error {
   name := filepath.Join(elem...)
   err := os.MkdirAll(filepath.Dir(name), os.ModeDir)
   if err != nil {
      return err
   }
   os.Stdout.WriteString("Create " + name + "\n")
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(value)
}

// Use 0 for INFO, 1 for VERBOSE and any other value for QUIET.
type LogLevel int

func (l LogLevel) Dump(req *http.Request) error {
   switch l {
   case 0:
      os.Stdout.WriteString(req.Method)
      os.Stdout.WriteString(" ")
      os.Stdout.WriteString(req.URL.String())
      os.Stdout.WriteString("\n")
   case 1:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return err
      }
      if IsBinary(buf) {
         quote := strconv.Quote(string(buf))
         os.Stdout.WriteString(quote)
      } else {
         os.Stdout.Write(buf)
      }
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         os.Stdout.WriteString("\n")
      }
   }
   return nil
}

// mimesniff.spec.whatwg.org#binary-data-byte
func IsBinary(buf []byte) bool {
   for _, b := range buf {
      switch {
      case b <= 0x08,
      b == 0x0B,
      0x0E <= b && b <= 0x1A,
      0x1C <= b && b <= 0x1F:
         return true
      }
   }
   return false
}

type Progress struct {
   *http.Response
   content int64
   part, partLength time.Time
}

func NewProgress(src *http.Response) *Progress {
   var pro Progress
   pro.Response = src
   pro.part = time.Now()
   pro.partLength = time.Now()
   return &pro
}

func (p *Progress) Read(buf []byte) (int, error) {
   since := time.Since(p.part)
   if since >= time.Second/2 {
      os.Stdout.WriteString(Percent(p.content, p.ContentLength))
      os.Stdout.WriteString("\t")
      os.Stdout.WriteString(LabelSize(p.content))
      os.Stdout.WriteString("\t")
      os.Stdout.WriteString(p.getRate())
      os.Stdout.WriteString("\n")
      p.part = p.part.Add(since)
   }
   // Callers should always process the n > 0 bytes returned before considering
   // the error err.
   read, err := p.Body.Read(buf)
   p.content += int64(read)
   return read, err
}

func (p Progress) getRate() string {
   rate := float64(p.content) / time.Since(p.partLength).Seconds()
   return LabelRate(rate)
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

func LabelNumber[T Number](value T) string {
   return Label(value, "", " K", " M", " B", " T")
}

func LabelRate[T Number](value T) string {
   return Label(value, " B/s", " kB/s", " MB/s", " GB/s", " TB/s")
}

func LabelSize[T Number](value T) string {
   return Label(value, " B", " kB", " MB", " GB", " TB")
}

func Percent[T Number](value, total T) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

type Number interface {
   float64 | int | int64 | uint64
}
