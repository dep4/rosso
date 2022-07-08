package mech

import (
   "github.com/89z/rosso/hls"
)

func three[T hls.Item]() {
   var b *T
   (*b).Ext()
}
