package main

import (
   "fmt"
   "net/url"
)

func main() {
   {
      val := url.Values{
         "one": {"two"},
      }
      fmt.Printf("%q\n", val.Get("one")) // "two"
   }
   {
      val := url.Values{
         "one": {"two"},
      }
      fmt.Printf("%q\n", val["one"]) // ["two"]
   }
   {
      val := url.Values{
         "one": {"two", "three"},
      }
      fmt.Printf("%q\n", val.Get("one")) // "two"
   }
   {
      val := url.Values{
         "one": {"two", "three"},
      }
      fmt.Printf("%q\n", val["one"]) // ["two" "three"]
   }
}
