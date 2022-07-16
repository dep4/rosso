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
   Log_Level int // this needs to work with flag.IntVar
   status int
   client http.Client
}

var Default_Client = Client{
   Log_Level: 1,
   client: http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   },
   status: http.StatusOK,
}

func (self Client) Do(req *http.Request) (*http.Response, error) {
   switch self.Log_Level {
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
      if !strconv.String(buf) {
         buf = strconv.AppendQuote(nil, string(buf))
      }
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         buf = append(buf, '\n')
      }
      os.Stderr.Write(buf)
   }
   res, err := self.client.Do(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != self.status {
      return nil, errors.New(res.Status)
   }
   return res, nil
}

func (self Client) Get(ref string) (*http.Response, error) {
   req, err := http.NewRequest("GET", ref, nil)
   if err != nil {
      return nil, err
   }
   return self.Do(req)
}

func (self Client) Level(level int) Client {
   self.Log_Level = level
   return self
}

func (self Client) Redirect(fn Redirect_Func) Client {
   self.client.CheckRedirect = nil
   return self
}

func (self Client) Status(status int) Client {
   self.status = status
   return self
}

func (self Client) Transport(tr *http.Transport) Client {
   self.client.Transport = tr
   return self
}

type Redirect_Func func(*http.Request, []*http.Request) error
