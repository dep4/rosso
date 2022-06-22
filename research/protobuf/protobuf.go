package main

import (
   "fmt"
)

type message map[string]encoder

type encoder interface {
   encode()
}

type encoded interface {
   uint32 | uint64
}

type slice[T encoded] []T

func (slice[T]) encode(){}

func add[T encoded](mes message, key string, val T) {
   switch out := mes[key].(type) {
   case nil:
      mes[key] = slice[T]{val}
   case slice[T]:
      mes[key] = append(out, val)
   default:
      fmt.Println("error", val)
   }
}

func main() {
   mes := make(message)
   add(mes, "one", uint64(1))
   add(mes, "one", uint64(2))
   add(mes, "one", uint64(3))
   add(mes, "one", uint32(4))
   fmt.Println(mes)
}
