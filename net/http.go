package net

import (
   "bufio"
   "fmt"
   "github.com/89z/format"
   "io"
   "net/http"
   "net/textproto"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func ReadRequest(src io.Reader) (*http.Request, error) {
   var req http.Request
   text := textproto.NewReader(bufio.NewReader(src))
   // .Method
   sMethodPath, err := text.ReadLine()
   if err != nil {
      return nil, err
   }
   // GET /fdfe/details?doc=com.instagram.android HTTP/1.1
   methodPath := strings.Fields(sMethodPath)
   if len(methodPath) != 3 {
      return nil, textproto.ProtocolError(sMethodPath)
   }
   req.Method = methodPath[0]
   // .URL
   addr, err := url.Parse(methodPath[1])
   if err != nil {
      return nil, err
   }
   req.URL = addr
   // .URL.Host
   head, err := text.ReadMIMEHeader()
   if err != nil {
      return nil, err
   }
   req.URL.Host = head.Get("Host")
   // .Header
   req.Header = http.Header(head)
   // .ContentLength
   sLength := head.Get("Content-Length")
   if sLength != "" {
      length, err := strconv.ParseInt(sLength, 10, 64)
      if err != nil {
         return nil, err
      }
      req.ContentLength = length
   }
   // .Body
   req.Body = io.NopCloser(text.R)
   return &req, nil
}

type Progress struct {
   *http.Response
   Content, part, PartLength int64
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
      meter := format.PercentInt64(p.Content, p.ContentLength)
      meter += "\t" + format.Size.LabelInt(p.Content)
      meter += "\t" + format.Rate.LabelInt(1000 * p.Content / end)
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
