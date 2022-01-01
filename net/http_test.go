package net

import (
   "os"
   "strings"
   "testing"
)

var sites = []string{
// keep this as content-length is required
`POST /api/v1/playlist/getFragment HTTP/1.1
Host: pandora.com
content-type:application/json
cookie:csrftoken=842b12c83a3c5153
x-authtoken:BXoTKywEhnoiEqDEcu0U/qGlFBEK5Tjblz3fgnLPgFojficRTR8Xm6Lw==
x-csrftoken:842b12c83a3c5153
content-length: 54

{"isStationStart":true,"stationId":126608766085892525}`,
`POST /player/addToPlaylist HTTP/1.1
Host: bleep.com
Content-Type: application/x-www-form-urlencoded

id=8728&type=ReleaseProduct`,
`GET /manifest.json HTTP/1.1
Host: github.com

`,
}

func TestRequest(t *testing.T) {
   for _, site := range sites {
      req, err := ReadRequest(strings.NewReader(site), true)
      if err != nil {
         t.Fatal(err)
      }
      if err := WriteRequest(os.Stdout, req); err != nil {
         t.Fatal(err)
      }
   }
}
