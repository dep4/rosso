package proto

import (
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

var data = []byte("Rw\b\x01\x13\x18\xb8\b \xf0\x10\x14*ihttps://play-lh.googleusercontent.com/mLw6yXfn-6EyvgvdYS7CQkUG0O7vLLQ4rtOvHf6Rq1_n5h9qwmB-Mtc1293CuLGkABAH\x01")

func TestProto(t *testing.T) {
   toks := parseUnknown(data)
   fmt.Printf("%+v\n", toks)
   var msg protopack.Message
   msg.UnmarshalAbductive(data, nil)
   fmt.Printf("%+v\n", msg)
}
