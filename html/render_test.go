package html

import (
   "os"
   "strings"
   "testing"
)

const s = `
<meta charset="utf-8">
<head>
<title>Umber</title>
<link rel="icon" href="/umber/media/umber.png">
</head>
`

func TestRender(t *testing.T) {
   l := NewLexer(strings.NewReader(s))
   err := l.Render(os.Stdout, " ")
   if err != nil {
      t.Fatal(err)
   }
}
