package protobuf

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "testing"
)

func TestGet(t *testing.T) {
   get(tag{1,"payload"})
}
