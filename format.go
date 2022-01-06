package format

import (
   "bytes"
   "io"
   "mime"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
   "strings"
   "time"
)

var Log = Logger{Writer: os.Stderr}

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
)

func Clean(char rune) rune {
   if strings.ContainsRune(`"*/:<>?\|`, char) {
      return -1
   }
   return char
}

// github.com/golang/go/issues/22318
func ExtensionByType(typ string) (string, error) {
   justType, _, err := mime.ParseMediaType(typ)
   if err != nil {
      return "", err
   }
   switch justType {
   case "audio/mp4":
      return ".m4a", nil
   case "audio/webm":
      return ".weba", nil
   case "video/mp4":
      return ".m4v", nil
   case "video/webm":
      return ".webm", nil
   }
   return "", notFound{justType}
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

func Trim(w io.Writer, s string) (int, error) {
   if len(s) >= 100 {
      s = s[:48] + "..." + s[len(s)-48:]
   }
   return io.WriteString(w, s)
}

type InvalidSlice struct {
   Index, Length int
}

func (i InvalidSlice) Error() string {
   index, length := int64(i.Index), int64(i.Length)
   var buf []byte
   buf = append(buf, "index out of range ["...)
   buf = strconv.AppendInt(buf, index, 10)
   buf = append(buf, "] with length "...)
   buf = strconv.AppendInt(buf, length, 10)
   return string(buf)
}

// Use 0 for INFO, 1 for VERBOSE and any other value for QUIET.
type Logger struct {
   io.Writer
   Level int
}

func (l Logger) Dump(req *http.Request) error {
   switch l.Level {
   case 0:
      s := req.Method + " " + req.URL.String() + "\n"
      io.WriteString(l, s)
   case 1:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return err
      }
      if IsBinary(buf) {
         s := strconv.Quote(string(buf))
         io.WriteString(l, s)
      } else {
         l.Write(buf)
      }
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         io.WriteString(l, "\n")
      }
   }
   return nil
}

type Progress struct {
   *http.Response
   content, partLength, part int64
   io.Writer
   time.Time
}

func NewProgress(src *http.Response, dst io.Writer) *Progress {
   var pro Progress
   pro.Response = src
   pro.Time = time.Now()
   pro.Writer = dst
   pro.partLength = 10_000_000
   return &pro
}

func (p *Progress) Read(buf []byte) (int, error) {
   if p.part == 0 {
      end := time.Since(p.Time).Milliseconds()
      if end >= 1 {
         PercentInt64(p.Writer, p.content, p.ContentLength)
         io.WriteString(p.Writer, "\t")
         Size.Int64(p.Writer, p.content)
         io.WriteString(p.Writer, "\t")
         Rate.Int64(p.Writer, 1000 * p.content / end)
         io.WriteString(p.Writer, "\n")
      }
   }
   // Callers should always process the n > 0 bytes returned before considering
   // the error err.
   read, err := p.Body.Read(buf)
   p.content += int64(read)
   p.part += int64(read)
   if p.part >= p.partLength {
      p.part = 0
   }
   return read, err
}

type notFound struct {
   input string
}

func (n notFound) Error() string {
   return strconv.Quote(n.input) + " not found"
}

type Symbols []string

func (s Symbols) Float64(w io.Writer, f float64) (int, error) {
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
   if i >= 1 {
      i = 3
   }
   symbol = strconv.FormatFloat(f, 'f', i, 64) + symbol
   return io.WriteString(w, symbol)
}

func (s Symbols) Int(w io.Writer, i int) (int, error) {
   f := float64(i)
   return s.Float64(w, f)
}

func (s Symbols) Int64(w io.Writer, i int64) (int, error) {
   f := float64(i)
   return s.Float64(w, f)
}

func (s Symbols) Uint64(w io.Writer, i uint64) (int, error) {
   f := float64(i)
   return s.Float64(w, f)
}

func Percent(w io.Writer, value, total float64) (int, error) {
   var s string
   if total != 0 {
      ratio := 100 * value / total
      s = strconv.FormatFloat(ratio, 'f', 1, 64)
   } else {
      s = "0"
   }
   return io.WriteString(w, s + "%")
}

func PercentInt(w io.Writer, value, total int) (int, error) {
   val, tot := float64(value), float64(total)
   return Percent(w, val, tot)
}

func PercentInt64(w io.Writer, value, total int64) (int, error) {
   val, tot := float64(value), float64(total)
   return Percent(w, val, tot)
}

func PercentUint64(w io.Writer, value, total uint64) (int, error) {
   val, tot := float64(value), float64(total)
   return Percent(w, val, tot)
}
