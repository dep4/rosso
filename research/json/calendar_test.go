package calendar

import (
   "bytes"
   "fmt"
   "testing"
)

func TestDate(t *testing.T) {
   in := date{Month: 12, Day: 31}
   buf := new(bytes.Buffer)
   in.WriteTo(buf)
   var out date
   out.ReadFrom(buf)
   fmt.Printf("%+v\n", out)
}
