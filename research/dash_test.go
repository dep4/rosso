package dash

import (
   "os"
   "testing"
)

const hexKey = "13d7c7cf295444944b627ef0ad2c1b3c"

func TestDASH(t *testing.T) {
   src, err := os.Open("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer src.Close()
   dst, err := os.Create("ignore/dec.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dst.Close()
   if err := start(src, dst, hexKey); err != nil {
      t.Fatal(err)
   }
}
