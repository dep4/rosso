package http

import (
   "bytes"
   "errors"
   "github.com/89z/rosso/strconv"
   "net/http"
   "net/http/httputil"
   "os"
)

type Client struct {
   client http.Client
   error_response bool
   level int
   status int
}

func (c Client) Error(response bool) Client {
   c.error_response = response
   return c
}

func (c Client) Redirect(fn Redirect_Func) Client {
   c.client.CheckRedirect = nil
   return c
}

func (c Client) Get(ref string) (*http.Response, error) {
   req, err := http.NewRequest("GET", ref, nil)
   if err != nil {
      return nil, err
   }
   return c.Do(req)
}

func (c Client) Status(status int) Client {
   c.status = status
   return c
}

func (c Client) Transport(tr *http.Transport) Client {
   c.client.Transport = tr
   return c
}

type Redirect_Func func(*http.Request, []*http.Request) error

var Default_Client = Client{
   client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
   status: http.StatusOK,
   error_response: true,
   level: info,
}

const (
   quiet = 0
   warning = 1
   info = 2
   verbose = 3
)

func (c Client) Quiet() Client {
   c.level = quiet
   return c
}

func (c Client) Warning() Client {
   c.level = warning
   return c
}

func (c Client) Info() Client {
   c.level = info
   return c
}

func (c Client) Verbose() Client {
   c.level = verbose
   return c
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
   switch c.level {
   case info:
      os.Stderr.WriteString(req.Method)
      os.Stderr.WriteString(" ")
      os.Stderr.WriteString(req.URL.String())
      os.Stderr.WriteString("\n")
   case verbose:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return nil, err
      }
      if !strconv.Valid(buf) {
         buf = strconv.AppendQuote(nil, buf)
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
      if c.error_response {
         return nil, errors.New(res.Status)
      } else if c.level >= warning {
         os.Stderr.WriteString(res.Status + "\n")
      }
   }
   return res, nil
}
