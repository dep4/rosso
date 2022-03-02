package hls

import (
)

type Segment struct {
   Key struct {
      method string
      uri string
   }
   Inf []struct {
      duration string
      uri string
   }
}
