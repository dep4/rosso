package main

import (
   "fmt"
)

type message map[int]interface{}

type token[T message | string]struct {
   in message
   out T
}

func main() {
   docV2 := message{
      13: message{
         1: message{
            34: message{2: "size"},
         },
      },
   }
   tok := token[message]{in: docV2}
   fmt.Println(tok)
}
