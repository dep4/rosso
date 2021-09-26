# September 25 2021

JSON method causes "b" output in some cases

For some reason, running this program:

~~~go
package main

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/js"
)

type visit struct{}

func (visit) Exit(js.INode) {}

func (v visit) Enter(n js.INode) js.IVisitor {
   n.JSON()
   return v
}

func main() {
   ast, err := js.Parse(parse.NewInputString("var x = {y: 5 > 2}"))
   if err != nil {
      panic(err)
   }
   var vis visit
   js.Walk(vis, ast)
}
~~~

Produces this output:

~~~
b
b
~~~

I am not sure why its outputting `b`, as I do not have any print statements in
the program.
