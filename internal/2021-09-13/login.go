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
   j.SortUsers()
   md5 := "37f691b063c10372135db21579643bf1"
   fmt.Println(md5)
   for _, user := range j.FilterUsers(md5) {
      fmt.Println(user.Count, user.Agent)
   }
}
