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
   Status_Code int
   http.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
   c.CheckRedirect = func(*http.Request, []*http.Request) error {
      return http.ErrUseLastResponse
   }
   c.Status_Code = http.StatusOK
   c.Transport = nil
   return c.Custom(req)
}

func (c Client) Custom(req *http.Request) (*http.Response, error) {
   switch c.Level {
   case 0:
      os.Stderr.WriteString(req.Method)
      os.Stderr.WriteString(" ")
      os.Stderr.WriteString(req.URL.String())
      os.Stderr.WriteString("\n")
   case 1:
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
   if res.StatusCode != c.Status_Code {
      return nil, errors.New(res.Status)
   }
   return res, nil
}
