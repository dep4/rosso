package net

import (
   "fmt"
   "testing"
)

const iframe = "<iframe  id='partnerPlayer' frameborder='0' marginwidth='0' marginheight='0' scrolling='no' width='100%' height='100%' src='//player.pbs.org/partnerplayer/wwGgFRSNeKGrsgjdYh6efQ==/?topbar=false&end=0&endscreen=true&start=0&autoplay=false' allowfullscreen></iframe>"

func TestURL(t *testing.T) {
   addr, err := ParseURL(iframe, " src='")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", addr)
}
