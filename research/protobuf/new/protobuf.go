package main

import (
   "encoding/json"
   "os"
)

type Message map[int]any

type String struct {
   Raw string
   Message
}

func main() {
   a := String{
      Raw: "hello",
      Message: Message{1: "world"},
   }
   json.NewEncoder(os.Stdout).Encode(a)
}
