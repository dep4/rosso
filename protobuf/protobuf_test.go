package protobuf

import (
   "fmt"
   "testing"
)

var checkin = Message{
   Tag{4, "checkin"}:Message{
      Tag{1, "build"}:Message{
         Tag{10, "sdkVersion"}: uint64(29),
      },
      Tag{2, ""}:Message{
         Tag{10, "sdkVersion"}: uint64(29),
      },
   },
}

func TestGet(t *testing.T) {
   {
      mes := checkin.Get(4, "checkin").Get(1, "build")
      fmt.Println(mes)
   }
   {
      mes := checkin.Get(1, "hello").Get(2, "world")
      fmt.Println(mes)
   }
}
