package xml

import (
   "bytes"
   "encoding/xml"
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
   _, after, found := bytes.Cut(buf, sep)
   if !found {
      return notFound(sep)
   }
   dec := xml.NewDecoder(bytes.NewReader(after))
   for {
      _, err := dec.Token()
      if err != nil {
         high := dec.InputOffset()
         return xml.Unmarshal(after[:high], val)
      }
   }
}

type notFound []byte

func (n notFound) Error() string {
   str := string(n)
   return strconv.Quote(str) + " is not found"
}
