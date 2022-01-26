package protobuf

import (
   "fmt"
   "testing"
)

func TestAdd(t *testing.T) {
   {
      mes := Message{
         {1, "hello"}: "world",
      }
      hello := mes.GetString(1, "hello")
      fmt.Printf("%q\n", hello)
   }
   {
      var mes Message
      err := mes.Add(1, "hello", Message{})
      fmt.Println(err)
   }
}
