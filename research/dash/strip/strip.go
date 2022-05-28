package main

import (
   "bytes"
   "encoding/hex"
   "os"
)

func main() {
   src, err := os.ReadFile("segment0.m4f")
   if err != nil {
      panic(err)
   }
   piff, err := hex.DecodeString("a2394f525a9b4f14a2446c427c648df4")
   if err != nil {
      panic(err)
   }
   free, err := hex.DecodeString("6672656500110010800000AA00389B71")
   if err != nil {
      panic(err)
   }
   pos := bytes.Index(src, piff)
   dst, err := os.Create("free.m4f")
   if err != nil {
      panic(err)
   }
   defer dst.Close()
   dst.Write(src[:pos])
   dst.Write(free)
   dst.Write(src[pos+len(piff):])
}
