package protobuf

import (
   "os"
   "testing"
)

func TestCheckin(t *testing.T) {
   for _, test := range tests {
      buf, err := os.ReadFile(test.name)
      if err != nil {
         t.Fatal(err)
      }
      mes, err := Unmarshal(buf)
      if err != nil {
         t.Fatal(err)
      }
      size := len(mes.Marshal())
      if size != test.size {
         t.Fatal(size)
      }
   }
}

type testType struct {
   name string
   size int
}

var tests = []testType {
   {"checkin.txt", 374},
   {"details.txt", 10095},
}
