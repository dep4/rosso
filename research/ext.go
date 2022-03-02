package main

import (
   "fmt"
   "mime"
)

func main() {
   err := mime.AddExtensionType(".webz", "audio/webz")
   if err != nil {
      panic(err)
   }
   ext, err := mime.ExtensionsByType("audio/webz")
   if err != nil {
      panic(err)
   }
   fmt.Println(ext)
}
