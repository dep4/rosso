package protobuf

import (
   "bufio"
   "bytes"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "os"
   "testing"
)

func TestNew(t *testing.T) {
   b := protowire.AppendBytes(nil, []byte("hello"))
   b = append(b, "world\n"...)
   r := bufio.NewReader(bytes.NewReader(b))
   v, err := consumeBytes(r)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(string(v))
   os.Stdout.ReadFrom(r)
}
