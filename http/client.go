package http

import (
   "bytes"
   "errors"
   "github.com/89z/format"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

type Client struct {
   http.Client
   Level int // this needs to work with flag.IntVar
   Status int
}

var Default_Client = Client{
   Client: http.Client{
      CheckRedirect: func(*Request, []*Request) error {
         return http.ErrUseLastResponse
      },
   },
   Level: 1,
   Status: http.StatusOK,
}

func (c Client) WithRedirect() Client {
   c.CheckRedirect = nil
   return c
}

func (c Client) WithLevel(level int) Client {
   c.Level = level
   return c
}

func (c Client) WithStatus(status int) Client {
   c.Status = status
   return c
}

func (c Client) WithTransport(tr *Transport) Client {
   c.Transport = tr
   return c
}

func (c Client) Do(req *Request) (*http.Response, error) {
   switch c.Level {
   case 1:
      os.Stderr.WriteString(req.Method)
      os.Stderr.WriteString(" ")
      os.Stderr.WriteString(req.URL.String())
      os.Stderr.WriteString("\n")
   case 2:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return nil, err
      }
      if !format.String(buf) {
         buf = strconv.AppendQuote(nil, string(buf))
      }
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         buf = append(buf, '\n')
      }
      os.Stderr.Write(buf)
   }
   res, err := c.Client.Do(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != c.Status {
      return nil, errors.New(res.Status)
   }
   return res, nil
}

func (c Client) Get(addr string) (*http.Response, error) {
   req, err := http.NewRequest("GET", addr, nil)
   if err != nil {
      return nil, err
   }
   return c.Do(req)
}
