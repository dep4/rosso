package dash

import (
   "fmt"
   "testing"
)

func Test_Reduce(t *testing.T) {
   rep := Representations{
      {ID: "one", Bandwidth: 1},
      {ID: "four", Bandwidth: 4},
      {ID: "seven", Bandwidth: 7},
      {ID: "ten", Bandwidth: 10},
      {ID: "thirteen", Bandwidth: 13},
   }.Filter(func(r Representation) bool {
      return r.Bandwidth > 2
   }).Reduce(func(carry, item Representation) bool {
      callback := Bandwidth(7)
      return callback(item) < callback(carry)
   })
   fmt.Printf("%+v\n", rep)
}
