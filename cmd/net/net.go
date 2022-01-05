package main

import (
   "bytes"
   "encoding/json"
   "github.com/89z/format/protobuf"
   "io"
   "os"
)

func readFrom(src io.Reader, output string, proto bool) error {
   dst, err := os.Create(output)
   if err != nil {
      dst = os.Stdout
   }
   defer dst.Close()
   if proto {
      mes, err := protobuf.Decode(src)
      if err != nil {
         return err
      }
      buf, err := json.Marshal(mes)
      if err != nil {
         return err
      }
      indent := new(bytes.Buffer)
      if err := json.Indent(indent, buf, "", " "); err != nil {
         return err
      }
      if _, err := dst.ReadFrom(indent); err != nil {
         return err
      }
   } else {
      _, err := dst.ReadFrom(src)
      if err != nil {
         return err
      }
   }
   return nil
}
