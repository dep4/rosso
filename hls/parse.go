package hls

import (
   "encoding/hex"
   "strings"
)

func hexDecode(s string) ([]byte, error) {
   s = strings.TrimPrefix(s, "0x")
   return hex.DecodeString(s)
}
