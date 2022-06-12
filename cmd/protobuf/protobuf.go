package main

import (
   "encoding/json"
   "github.com/89z/format/protobuf"
   "os"
)

func doProtoBuf(input, output string) error {
   data, err := os.ReadFile(input)
   if err != nil {
      return err
   }
   mes := make(protobuf.Message)
   if err := mes.UnmarshalBinary(data); err != nil {
      return err
   }
   file, err := os.Create(output)
   if err != nil {
      file = os.Stdout
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", " ")
   return enc.Encode(mes)
}
