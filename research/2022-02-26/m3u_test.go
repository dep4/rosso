package m3u

import (
   "fmt"
   "testing"
)

var tests = []string{
   `D:\two\Foo.mp3`,
   `http://example.com/two/Mine.mp3`,
   `two/Stuff.mp3`,
   `two\Stuff.mp3`,
}

func TestM3U(t *testing.T) {
   for _, test := range tests {
      spec := specification(test)
      fmt.Println(spec)
   }
}
