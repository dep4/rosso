package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "os"
)

const target = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36"

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
   for _, user := range j.Users {
      if user.Agent == target {
         fmt.Println(user)
         fmt.Println(j.JA3(user.MD5))
         fmt.Println()
      }
   }
}

var sites = []string{
   "https://www.reddit.com",
   "https://github.com",
   "https://stackoverflow.com",
   "https://variety.com",
   "https://vimeo.com",
   "https://www.google.com",
   "https://www.indiewire.com",
   "https://www.nytimes.com",
   "https://www.wikipedia.org",
   "https://www.youtube.com",
}

func version(min uint16) []uint16 {
   vs := []uint16{772, 771, 770, 769, 768}
   for k, v := range vs {
      if v == min {
         return vs[:k+1]
      }
   }
   return nil
}
