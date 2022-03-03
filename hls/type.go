package hls

import (
   "crypto/cipher"
)

type Master struct {
   Media []Media
   Stream []Stream
}

type Media struct {
   GroupID string
   URI string
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle audio only
   Audio string // link to Media
   URI string
}

type Decrypter struct {
   cipher.Block
   IV []byte
}

type Information struct {
   Duration string
   URI string
}

type Segment struct {
   Key struct {
      Method string
      URI string
   }
   Info []Information
}
