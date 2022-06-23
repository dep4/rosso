package format

import (
   "bytes"
   "errors"
   "net/http"
   "net/http/httputil"
   "os"
   "strconv"
)

func check_redirect(*http.Request, []*http.Request) error {
   return http.ErrUseLastResponse
}

type Client struct {
   Log_Level int // this needs to work with flag.IntVar
   http.Client
}

func New_Client() Client {
   var c Client
   c.CheckRedirect = check_redirect
   return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
   switch c.Log_Level {
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
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   c.CheckRedirect = check_redirect
   c.Transport = nil
   return res, nil
}
