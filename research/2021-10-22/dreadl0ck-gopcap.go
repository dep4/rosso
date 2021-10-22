package main

import (
   "fmt"
   "github.com/dreadl0ck/gopcap"
)

func main() {
   f, err := gopcap.Open("PCAPdroid_22_Oct_15_19_28.pcap")
   if err != nil {
      panic(err)
   }
   defer f.Close()
   for {
      head, body, err := f.ReadNextPacket()
      if err != nil {
         break
      }
      fmt.Printf("%+v\n", head)
      fmt.Printf("%q\n\n", body)
   }
}
