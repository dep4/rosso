package tool

import (
	"fmt"
	"io"
	"net/http"
	"time"
      "github.com/89z/format"
)

func Get(url string) (io.ReadCloser, error) {
   c := http.Client{
   Timeout: time.Duration(60) * time.Second,
   }
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
   return nil, err
   }
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
   format.LogLevel.Dump(1, req)
   resp, err := c.Do(req)
   if err != nil {
   return nil, err
   }
   if resp.StatusCode != 200 {
   return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
   }
   return resp.Body, nil
}
