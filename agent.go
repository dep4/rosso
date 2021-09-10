package ja3

import (
   "encoding/json"
   "io"
)

const AllUas = "https://ja3er.com/getAllUasJson"

type Agent struct {
   LastSeen string `json:"Last_seen"`
   MD5 string
   UserAgent string `json:"User-Agent"`
}

func Agents(r io.Reader) ([]Agent, error) {
   var a []Agent
   if err := json.NewDecoder(r).Decode(&a); err != nil {
      return nil, err
   }
   return a, nil
}
