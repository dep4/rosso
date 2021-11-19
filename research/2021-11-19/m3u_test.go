package m3u

import (
   "fmt"
   "testing"
)

var stream = []byte(
   "#EXT-X-STREAM-INF:" +
   `CLOSED-CAPTIONS=NONE,BANDWIDTH=1056000,AVERAGE-BANDWIDTH=1051000,` +
   `RESOLUTION=218x432,CODECS="avc1.4d001e,mp4a.40.2",AUDIO="1"`,
)

func TestPlaylist(t *testing.T) {
   buf := newBuffer(stream)
   buf.readBytes(':', '"')
   for {
      field := buf.readBytes(',', '"')
      if field == nil {
         break
      }
      fmt.Printf("%s\n", field)
   }
}
