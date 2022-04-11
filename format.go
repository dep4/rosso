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

func Create[T any](value T, elem ...string) error {
   name := filepath.Join(elem...)
   err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
   if err != nil {
      return err
   }
   os.Stderr.WriteString("Create " + name + "\n")
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(value)
}

func Open[T any](elem ...string) (*T, error) {
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
      p.progress()
      p.part = p.part.Add(since)
   }
   // Callers should always process the n > 0 bytes returned before considering
   // the error err.
   read, err := p.Body.Read(buf)
   p.content += int64(read)
   return read, err
}

func (p Progress) progress() {
   rate := float64(p.content) / time.Since(p.partLength).Seconds()
   os.Stderr.WriteString(Percent(p.content, p.ContentLength))
   os.Stderr.WriteString("\t")
   os.Stderr.WriteString(LabelSize(p.content))
   os.Stderr.WriteString("\t")
   os.Stderr.WriteString(LabelRate(rate))
   os.Stderr.WriteString("\n")
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

func Percent[T, U Number](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

type Number interface {
   float64 | int | int64 | ~uint64
}

type LogLevel int

func (l LogLevel) Dump(req *http.Request) error {
   quote := func(b []byte) []byte {
      if IsBinary(b) {
         b = strconv.AppendQuote(nil, string(b))
      }
      if !bytes.HasSuffix(b, []byte{'\n'}) {
         b = append(b, '\n')
      }
      return b
   }
   switch l {
   case 0:
      os.Stderr.WriteString(req.Method)
      os.Stderr.WriteString(" ")
      os.Stderr.WriteString(req.URL.String())
      os.Stderr.WriteString("\n")
   case 1:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return err
      }
      os.Stderr.Write(quote(buf))
   case 2:
      buf, err := httputil.DumpRequestOut(req, true)
      if err != nil {
         return err
      }
      os.Stderr.Write(quote(buf))
   }
   return nil
}
