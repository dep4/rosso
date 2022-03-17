package main

import (
   "fmt"
)

type message map[int]interface{}

func main() {
   docV2 := message{
      13: message{
         1: message{
            34: message{2: "size"},
         },
      },
   }
   fmt.Println(docV2)
}
