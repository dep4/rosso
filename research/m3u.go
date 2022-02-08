package m3u

import (
   "time"
)

type Master struct {
   Version int
   Stream []Stream
}

type Stream struct {
   Codecs string
   URI string
}
