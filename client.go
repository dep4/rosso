package format

import (
   "bytes"
   "errors"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

type Client struct {
   Level int // this needs to work with flag.IntVar
   Status int
   http.Client
}

var Default_Client = Client{
   Level: 1,
   Status: http.StatusOK,
   Client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
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
      if !String(buf) {
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

func (c Client) WithLevel(level int) Client {
   c.Level = level
   return c
}

func (c Client) WithRedirect() Client {
   c.CheckRedirect = nil
   return c
}

func (c Client) WithStatus(status int) Client {
   c.Status = status
   return c
}

func (c Client) WithTransport(tr *http.Transport) Client {
   c.Transport = tr
   return c
}
