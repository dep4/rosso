package main

import (
   "encoding/csv"
   "fmt"
   "strings"
)

const media = `TYPE=SUBTITLES,GROUP-ID="subtitles_vod-ap3-aoc.tv.apple.com",LANGUAGE="en",NAME="English",AUTOSELECT=YES,FORCED=NO,STABLE-RENDITION-ID="ffa9aea0ca66b8731fe54339d311b8b234d0a3becffb8c3312adf77df048c146",URI="stream/playlist.m3u8?cc=SI&cdn=vod-ap3-aoc.tv.apple.com&a=1484589502&p=461374806&st=1821682691&a=1524726231&p=377684155&st=1491017020&a=1524197777&p=368330428&st=1449517832&a=1524197722&p=368330432&st=1449518338&a=1524198082&p=368330370&st=1449517662&a=1525078430&p=368329706&st=1449510111&a=1524197604&p=368330236&st=1449518904&a=1524197554&p=368330322&st=1449518550&a=1524197773&p=368330253&st=1449518005&a=1539152595&p=368283705&st=1449199987"`

var tests = []string{
   `ignore\apple-master.m3u8`,
   `m3u8\cbc-master.m3u8`,
   `m3u8\nbc-master.m3u8`,
   `m3u8\paramount-master.m3u8`,
   `m3u8\roku-master.m3u8`,
}

func main() {
   r := csv.NewReader(strings.NewReader(media))
   r.LazyQuotes = true
   record, err := r.Read()
   if err != nil {
      panic(err)
   }
   for _, field := range record {
      fmt.Printf("%q\n", field)
   }
}
