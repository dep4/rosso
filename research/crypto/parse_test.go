package crypto

import (
   "fmt"
   "testing"
   "github.com/wangluozhe/requests"
   "github.com/wangluozhe/requests/url"
)

// error decoding message
const hello = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513-21,29-23-24,0"

func TestTwo(t *testing.T) {
   req := url.NewRequest()
   req.Ja3 = hello
   r, err := requests.Get("https://android.googleapis.com/auth", req)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(r.Text)
}
