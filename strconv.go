package format

import (
   "bytes"
   "fmt"
   "mime"
   "net/http"
   "net/http/httputil"
   "strconv"
   "strings"
   "time"
)

var (
   Number = Symbols{"", " K", " M", " B", " T"}
   Size = Symbols{" B", " kB", " MB", " GB", " TB"}
   Rate = Symbols{" B/s", " kB/s", " MB/s", " GB/s", " TB/s"}
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

func (p Progress) Range() string {
   buf := []byte("bytes=")
   buf = strconv.AppendInt(buf, p.Content, 10)
   buf = append(buf, '-')
   buf = strconv.AppendInt(buf, p.Content+p.PartLength-1, 10)
   return string(buf)
}

type Symbols []string

type notFound struct {
   input string
}

func (n notFound) Error() string {
   return strconv.Quote(n.input) + " not found"
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

type Progress struct {
   *http.Response
   Content, PartLength, part int64
   time.Time
}

func Content(length int64) Progress {
   pro := Progress{PartLength: 10_000_000}
   pro.Response = &http.Response{ContentLength: length}
   pro.Time = time.Now()
   return pro
}

func Response(res *http.Response) *Progress {
   pro := Progress{Response: res, PartLength: 10_000_000}
   pro.Time = time.Now()
   return &pro
}

func Percent(value, total float64) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * value / total
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
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

func PercentInt(value, total int) string {
   val, tot := float64(value), float64(total)
   return Percent(val, tot)
}

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

// Use 0 for INFO, 1 for VERBOSE and any other value for QUIET.
type LogLevel int

func (l LogLevel) Dump(req *http.Request) error {
   switch l {
   case 0:
      loc := req.URL.String()
      loc = Trim(loc)
      fmt.Println(req.Method, loc)
   case 1:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return err
      }
      if IsBinary(buf) {
         fmt.Printf("%q", buf)
      } else {
         fmt.Printf("%s", buf)
      }
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         fmt.Println()
      }
   }
   return nil
}

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

func (s Symbols) LabelInt(i int64) string {
   f := float64(i)
   return s.Label(f)
}

func (s Symbols) LabelUint(i uint64) string {
   f := float64(i)
   return s.Label(f)
}
