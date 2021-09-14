package agent

import (
   "bytes"
   "fmt"
   "github.com/89z/parse/ja3"
   "io"
   "net/http"
   "time"
)

func readAll(r io.ReadCloser) ([]byte, error) {
   if r == nil {
      return nil, nil
   }
   return io.ReadAll(r)
}

func find(j *ja3.JA3er, req *http.Request) (*ja3.User, error) {
   data, err := readAll(req.Body)
   if err != nil {
      return nil, err
   }
   done := make(map[string]bool)
   for _, user := range j.Users {
      hello := j.JA3(user.MD5)
      if done[hello] {
         continue
      } else {
         done[hello] = true
      }
      spec, err := ja3.Parse(hello)
      if err != nil {
         fmt.Println(err)
         continue
      }
      req.Body = io.NopCloser(bytes.NewReader(data))
      time.Sleep(100 * time.Millisecond)
      res, err := ja3.NewTransport(spec).RoundTrip(req)
      if err != nil {
         fmt.Println(err)
         continue
      }
      defer res.Body.Close()
      fmt.Println(res.Status)
      if res.StatusCode == http.StatusOK {
         return &user, nil
      }
   }
   return nil, fmt.Errorf("%+v FAIL", req)
}
