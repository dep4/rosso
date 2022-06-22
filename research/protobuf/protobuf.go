package main

import (
   "fmt"
)

func add[T any](mes message, num string, val T) {
   in := mes[num]
   switch out := in.(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = slice[T]{out, val}
   case slice[T]:
      mes[num] = append(out, val)
   }
}

type message map[string]any

func (m message) add(num string, val any) {
   add(m, num, val)
}

type slice[T any] []T

func main() {
   checkin := make(message)
   for n := 0; n < 9; n++ {
      checkin.add("nine", 9)
   }
   fmt.Println(checkin)
}
