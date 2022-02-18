package googleplay

import (
   "encoding/json"
   "os"
   "testing"
)

func TestRequest(t *testing.T) {
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   req := newRequest()
   enc.Encode(req)
   enc.Encode(req.message())
}
