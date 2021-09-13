package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "os"
)

func main() {
   ua, err := os.Open("getAllUasJson.json")
   if err != nil {
      panic(err)
   }
   defer ua.Close()
   hash, err := os.Open("getAllHashesJson.json")
   if err != nil {
      panic(err)
   }
   defer hash.Close()
   j, err := ja3.NewJA3er(ua, hash)
   if err != nil {
      panic(err)
   }
   for _, agent := range j.Agents("37f691b063c10372135db21579643bf1") {
      fmt.Println(agent)
   }
}
