package protobuf

import (
   "fmt"
   "testing"
)

var (
   fail = []byte("\nU\xe2\x01R\nPCjkaNwoTNDU2NTczMDgzMDY3NDQ0Njc4MhIgChAxNjM2MDgzODEzNDE2MDUwEgwI5dCSjAYQ0NaxxgE=*\x03\b\xb7\x01J\x1c\b\x12\x9a\x01\x17\n\x13\b\xb4Ԍ\x94\xa7\x80\xf4\x02\x15\xb3\x83\x9a\x00\x1d\x19\xc3\x02\xe4\x10\x01")
   pass = []byte("\nU\xe2\x01R\nPCjkaNwoTNDM4MDc5NTIzOTc1MTU3Nzg3NxIgChAxNjM2MDgzNzYwOTgwODc4EgwIsNCSjAYQsIXc0wM=*\x03\b\x8b\x02J\x1c\b\x12\x9a\x01\x17\n\x13\b\x86\x89\x88\xfb\xa6\x80\xf4\x02\x15.\x85\x9a\x00\x1d\x93\x9f\b2\x10\x01")
   serverLogsCookie = []byte("\b\x12\x9a\x01\x17\n\x13\b\x86\x89\x88\xfb\xa6\x80\xf4\x02\x15.\x85\x9a\x00\x1d\x93\x9f\b2\x10\x01")
)

func TestProto(t *testing.T) {
   fields := Parse(pass)
   fmt.Print(fields, "\n\n")
   fields = Parse(fail)
   fmt.Print(fields, "\n\n")
   fields = Parse(serverLogsCookie)
   fmt.Println(fields)
}
