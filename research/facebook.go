package main

import (
   "fmt"
   "github.com/89z/format/json"
   "os"
)

func main() {
   file, err := os.Open("ignore.html")
   if err != nil {
      panic(err)
   }
   defer file.Close()
   scan, err := json.NewScanner(file)
   if err != nil {
      panic(err)
   }
   scan.Split = []byte(`{"\u0040context"`)
   scan.Scan()
   var media struct {
      DateCreated string
   }
   if err := scan.Decode(&media); err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", media)
   // "image":{"uri":"https:\/\/scontent-hou1-1.xx.fbcdn.net\/v\/t15.5256-10\/75
   // "text":"Who's Gonna Save My Soul | Gnarls Barkley | From The Basement"
}
