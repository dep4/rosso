package main

import (
   "encoding/base64"
)

var pssh = []byte{
   0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
   // Widevine UUID:
   0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
   0, 0, 0, 0,
   // length + KID:
   8, 0, 0, 0, 0, 0, 0, 0,
}

func main() {
   base64.StdEncoding.EncodeToString(pssh)
}
