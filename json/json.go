package json

import (
   "bytes"
   "encoding/json"
   "strconv"
)

func Unmarshal(buf, sep []byte, val any) error {
   _, after, found := bytes.Cut(buf, sep)
   if !found {
      return notFound(sep)
   }
   dec := json.NewDecoder(bytes.NewReader(after))
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return json.Unmarshal(after[:high], val)
      }
   }
}

type notFound []byte

func (n notFound) Error() string {
   str := string(n)
   return strconv.Quote(str) + " is not found"
}
