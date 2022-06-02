package mp4

import (
	"encoding/hex"
	"io"
	"os"
	"testing"
)

func TestDecrypt(t *testing.T) {
	for _, test := range tests {
		dst, err := os.Create(test.create)
		if err != nil {
			t.Fatal(err)
		}
		if err := writeInit(test.init, dst); err != nil {
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
		if err := DecryptMP4withCenc(seg, key, dst); err != nil {
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

type testType struct {
	init    string
	segment string
	create  string
	key     string
}

var tests = []testType{
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

func writeInit(src string, dst io.Writer) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}

