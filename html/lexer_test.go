package html

import (
   "fmt"
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

func TestNext(t *testing.T) {
   l := NewLexer(strings.NewReader(s))
   l.NextAttr("rel", "icon")
   fmt.Println(l.GetAttr("href"))
}
