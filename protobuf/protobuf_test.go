package protobuf

import (
   "os"
   "testing"
)

var buf = []byte("\nU\xe2\x01R\nPCjkaNwoTMzk0NzAwOTc1MjU3MDUyNzQ0MBIgChAxNjM2MDYyNjA5NzAzMjkxEgwIkauRjAYQ-LytzwI=*\x03\b\xf3\x01J\x1c\b\x12\x9a\x01\x17\n\x13\b\x84\xb3\xad\x95\xd8\xff\xf3\x02\x15a\x05\xc9\n\x1dz7\x05\xd4\x10\x01")

func TestProto(t *testing.T) {
   fields := Parse(buf)
   buf, err := Indent(fields)
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(buf)
}
