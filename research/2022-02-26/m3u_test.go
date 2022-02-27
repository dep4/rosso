package m3u

import (
   "testing"
)

type testType struct {
   m3u, mp3, res string
}

var tests = []testType{
   {`D:\one\Foo.m3u`,
   `D:\two\Foo.mp3`,
   `d:\two\Foo.mp3`},
   {`one\Stuff.m3u`,
   `D:\two\Foo.mp3`,
   `d:\two\Foo.mp3`},
   {`http://example.com/one/Mine.m3u`,
   `D:\two\Foo.mp3`,
   `d:\two\Foo.mp3`},
   {`one\Stuff.m3u`,
   `http://example.com/two/Mine.mp3`,
   `http://example.com/two/Mine.mp3`},
   {`D:\one\Foo.m3u`,
   `http://example.com/two/Mine.mp3`,
   `http://example.com/two/Mine.mp3`},
   {`http://example.com/one/Mine.m3u`,
   `http://example.com/two/Mine.mp3`,
   `http://example.com/two/Mine.mp3`},
   {`http://example.com/one/Mine.m3u`,
   `two/Stuff.mp3`,
   `http://example.com/one/two/Stuff.mp3`},
   {`D:/one/Foo.m3u`,
   `two/Stuff.mp3`,
   `D:\one\two\Stuff.mp3`},
   {`one/Stuff.m3u`,
   `two/Stuff.mp3`,
   `one/two/Stuff.mp3`},
}

func TestM3U(t *testing.T) {
   for _, test := range tests {
      //res := abs(test.m3u, test.mp3)
      res := resolve2(test.m3u, test.mp3)
      if res != test.res {
         t.Error(test, res)
      }
   }
}
