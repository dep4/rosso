package main

import (
   "encoding/hex"
   "fmt"
)

func main() {
   piff, err := hex.DecodeString("a2394f525a9b4f14a2446c427c648df4")
   if err != nil {
      panic(err)
   }
   fmt.Printf("%#v\n", piff)
}
