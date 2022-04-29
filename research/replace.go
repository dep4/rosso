package main

import (
   "fmt"
   "strconv"
   "strings"
)

func replace(s, id string, start int) string {
   s = strings.Replace(s, "$RepresentationID$", id, 1)
   s = strings.Replace(s, "$Time$", strconv.Itoa(start), 1)
   return s
}

const media =
   "CH4_44_7_900_18926001001003_001_J01-$RepresentationID$-$Time$.dash"

func main() {
   fmt.Println(replace(media, "hello", 9))
}
