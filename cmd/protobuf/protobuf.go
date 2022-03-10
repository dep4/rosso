package main

import (
   "encoding/json"
   "github.com/89z/format/protobuf"
   "os"
)

func doProtoBuf(input, output string) error {
   buf, err := os.ReadFile(input)
   if err != nil {
      return err
   }
   mes, err := protobuf.Unmarshal(buf)
   if err != nil {
      return err
   }
   dst, err := os.Create(output)
   if err != nil {
      dst = os.Stdout
   }
   defer dst.Close()
   enc := json.NewEncoder(dst)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc.Encode(mes)
}
