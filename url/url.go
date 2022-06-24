package url

import (
   "bufio"
   "io"
   "net/url"
   "strings"
)

func Decode(r io.Reader) (url.Values, error) {
   vals := make(url.Values)
   buf := bufio.NewReader(r)
   for {
      key, err := buf.ReadString('=')
      if err == io.EOF {
         break
      } else if err != nil {
         return nil, err
      }
      val, err := buf.ReadString('\n')
      key = strings.TrimSuffix(key, "=")
      val = strings.TrimSuffix(val, "\n")
      vals.Add(key, val)
      if err == io.EOF {
         break
      } else if err != nil {
         return nil, err
      }
   }
   return vals, nil
}

func Encode(w io.Writer, vals url.Values) error {
   for key := range vals {
      val := vals.Get(key)
      if _, err := io.WriteString(w, key); err != nil {
         return err
      }
      if _, err := io.WriteString(w, "="); err != nil {
         return err
      }
      if _, err := io.WriteString(w, val); err != nil {
         return err
      }
      if _, err := io.WriteString(w, "\n"); err != nil {
         return err
      }
   }
   return nil
}
