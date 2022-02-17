package net

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

// keep this as content-length is required
const post = `POST /api/v1/playlist/getFragment HTTP/1.1
Host: pandora.com
content-type:application/json
cookie:csrftoken=842b12c83a3c5153
x-authtoken:BXoTKywEhnoiEqDEcu0U/qGlFBEK5Tjblz3fgnLPgFojficRTR8Xm6Lw==
x-csrftoken:842b12c83a3c5153
content-length: 54

{"isStationStart":true,"stationId":126608766085892525}`

func TestRequest(t *testing.T) {
   req, err := ReadRequest(strings.NewReader(post), true)
   if err != nil {
      t.Fatal(err)
   }
   if err := WriteRequest(req, os.Stdout); err != nil {
      t.Fatal(err)
   }
}
