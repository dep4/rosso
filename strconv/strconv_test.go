package strconv

import (
   "testing"
)

func Test_Append(t *testing.T) {
   var b []byte
   b = AppendCardinal(nil, 123)
   if s := string(b); s != "123" {
      t.Fatal(s)
   }
   b = AppendCardinal(nil, 1234)
   if s := string(b); s != "1.23 thousand" {
      t.Fatal(s)
   }
   b = AppendSize(nil, 123)
   if s := string(b); s != "123 byte" {
      t.Fatal(s)
   }
   b = AppendSize(nil, 1234)
   if s := string(b); s != "1.23 kilobyte" {
      t.Fatal(s)
   }
   b = NewRatio(1234, 10).AppendRate(nil)
   if s := string(b); s != "123 byte/s" {
      t.Fatal(s)
   }
   b = NewRatio(12345, 10).AppendRate(nil)
   if s := string(b); s != "1.23 kilobyte/s" {
      t.Fatal(s)
   }
   b = NewRatio(1234, 10000).AppendPercent(nil)
   if s := string(b); s != "12.34%" {
      t.Fatal(s)
   }
}
