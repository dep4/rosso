package scanner

import (
   "fmt"
   "strings"
   "testing"
)

var src = strings.NewReader(`
#EXT-X-STREAM-INF:BANDWIDTH=400000,CODECS="mp4a.40.5,avc1.42C015",RESOLUTION=384x216,FRAME-RATE=25,AUDIO="audio-aach-96",CLOSED-CAPTIONS=NONE
`)

func TestScan(t *testing.T) {
   form := newFormat(src)
   fmt.Printf("%+v\n", form)
}
