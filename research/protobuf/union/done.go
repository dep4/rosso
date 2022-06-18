package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[protowire.Number]Token

type Token struct {
   Type protowire.Type
   Value []Value
}

type Value struct {
   Int32 uint32
   Int64 uint64
   Bytes []byte
   Message Message
}
