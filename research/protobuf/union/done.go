package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[protowire.Number]Token

type Token struct {
   Wire protowire.Type
   Value []Value
}

type Value struct {
   Integer uint64
   Bytes []byte
   Message Message
}
