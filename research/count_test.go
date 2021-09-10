package sort

import (
   "fmt"
   "os"
   "testing"
)

func TestCount(t *testing.T) {
   f, err := os.Open("getAllUasJson.json")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   cs, err := newCounts(f)
   if err != nil {
      t.Fatal(err)
   }
   g := cs.filter("2021-08-10").groups().filter(9)
   fmt.Println(len(g))
}
