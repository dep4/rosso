package json

import (
   "bytes"
   "encoding/json"
   "io"
   "strconv"
)

func Decode(src io.Reader, sep []byte, val any) error {
   buf, err := io.ReadAll(src)
   if err != nil {
      return err
   }
   return Unmarshal(buf, sep, val)
}

func Unmarshal(buf, sep []byte, val any) error {
   _, after, ok := bytes.Cut(buf, sep)
   if !ok {
      return notPresent(sep)
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

type notPresent []byte

func (n notPresent) Error() string {
   str := string(n)
   return strconv.Quote(str) + " is not present"
}
