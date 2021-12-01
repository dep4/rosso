package parse

import (
   "strconv"
)

type Invalid struct {
   Input string
}

func (i Invalid) Error() string {
   return strconv.Quote(i.Input) + " invalid"
}
