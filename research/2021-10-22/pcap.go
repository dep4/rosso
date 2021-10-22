package main

import (
   "github.com/dreadl0ck/ja3"
   "os"
)

/*
github.com/microsoft/etl2pcapng
*/
func main() {
   ja3.ReadFileJSON("PktMon.etl", os.Stdout, false)
}
