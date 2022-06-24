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

func (c Client) Redirect() Client {
   c.client.CheckRedirect = nil
   return c
}

func (c Client) Level(level int) Client {
   c.level = level
   return c
}

func (c Client) Status(status int) Client {
   c.status = status
   return c
}

func (c Client) Transport(tr *http.Transport) Client {
   c.client.Transport = tr
   return c
}

type Client struct {
   level int // this needs to work with flag.IntVar
   status int
   client http.Client
}

var Default_Client = Client{
   level: 1,
   status: http.StatusOK,
   client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
   switch c.level {
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
   res, err := c.client.Do(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != c.status {
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
