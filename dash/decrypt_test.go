package dash

import (
   "encoding/hex"
   "os"
   "testing"
)

type decryptTest struct {
	init    string
	segment string
	create  string
	key     string
}

var decTests = []decryptTest{
	{
		"ignore/amc-init0.m4f",
		"ignore/amc-segment0.m4f",
		"ignore/amc.mp4",
		"a66a5603545ad206c1a78e160a6710b1",
	},
	{
		"ignore/paramount-init.m4v",
		"ignore/paramount-seg_1.m4s",
		"ignore/paramount.mp4",
		"44f12639c9c4a5a432338aca92e38920",
	},
	{
		"ignore/roku-index_video_1_0_init.mp4",
		"ignore/roku-index_video_1_0_1.mp4",
		"ignore/roku.mp4",
		"13d7c7cf295444944b627ef0ad2c1b3c",
	},
}

func TestDecrypt(t *testing.T) {
   for _, test := range decTests {
      dst, err := os.Create(test.create)
      if err != nil {
         t.Fatal(err)
      }
      init, err := os.Open(test.init)
      if err != nil {
         t.Fatal(err)
      }
      dst.ReadFrom(init)
      if err := init.Close(); err != nil {
         t.Fatal(err)
      }
      seg, err := os.Open(test.segment)
      if err != nil {
         t.Fatal(err)
      }
      key, err := hex.DecodeString(test.key)
      if err != nil {
         t.Fatal(err)
      }
      if err := Decrypt(dst, seg, key); err != nil {
         t.Fatal(err)
      }
      if err := seg.Close(); err != nil {
         t.Fatal(err)
      }
      if err := dst.Close(); err != nil {
         t.Fatal(err)
      }
   }
}
