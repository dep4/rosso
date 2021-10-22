package main

import (
   "github.com/dreadl0ck/ja3"
   "os"
)

func main() {
   ja3.ReadFileJSON("PCAPdroid_22_Oct_15_19_28.pcap", os.Stdout, false)
}
