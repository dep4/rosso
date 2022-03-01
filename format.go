package format

import (
   "bytes"
   "mime"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
   "strings"
   "time"
)

func Clean(char rune) rune {
   if strings.ContainsRune(`"*/:<>?\|`, char) {
      return -1
   }
   return char
}

// github.com/golang/go/issues/22318
func ExtensionByType(typ string) (string, error) {
   mediaType, _, err := mime.ParseMediaType(typ)
   if err != nil {
      return "", err
   }
   var ext string
   switch mediaType {
   case "audio/mp4":
      ext = ".m4a"
   case "audio/mpeg":
      ext = ".mp3"
   case "audio/webm":
      ext = ".weba"
   case "video/mp4":
      ext = ".m4v"
   case "video/webm":
      ext = ".webm"
   }
   return ext, nil
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

func PercentInt(value, total int) string {
   val, tot := float64(value), float64(total)
   return percent(val, tot)
}

func PercentInt64(value, total int64) string {
   val, tot := float64(value), float64(total)
   return percent(val, tot)
}

func percent(value, total float64) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
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
      os.Stdout.WriteString(PercentInt64(p.content, p.ContentLength))
      os.Stdout.WriteString("\t")
      os.Stdout.WriteString(Size.GetInt64(p.content))
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
   return Rate.Get(rate)
}

type Symbols []string

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
)

func (s Symbols) Get(f float64) string {
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
   return strconv.FormatFloat(f, 'f', i, 64) + symbol
}

func (s Symbols) GetInt64(i int64) string {
   f := float64(i)
   return s.Get(f)
}

func (s Symbols) GetUint64(i uint64) string {
   f := float64(i)
   return s.Get(f)
}
