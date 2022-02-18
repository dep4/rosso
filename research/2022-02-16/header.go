package main

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "os"
)

func main() {
   buf, err := os.ReadFile("ignore.txt")
   if err != nil {
      panic(err)
   }
   responseWrapper, err := protobuf.Unmarshal(buf)
   if err != nil {
      panic(err)
   }
   deliveryData := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      Get(2, "appDeliveryData")
   a := deliveryData.Get(15, "splitDeliveryData")
   fmt.Printf("%+v\n", a)
   b := deliveryData.GetMessages(15, "splitDeliveryData")
   fmt.Printf("%+v\n", b)
}
