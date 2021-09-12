package ja3

import (
   "net/http"
   "testing"
)

const test =
   "771," +
   "49196-49195-49200-49199-159-158-49188-49187-49192-49191-49162-49161-49172-49171-157-156-61-60-53-47-10," +
   "0-10-11-13-35-23-65281,,"

func TestParse(t *testing.T) {
   spec, err := Parse(test)
   if err != nil {
      t.Fatal(err)
   }
   // https://www.reddit.com
   req, err := http.NewRequest("GET", "https://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   if _, err := NewTransport(spec).RoundTrip(req); err != nil {
      t.Fatal(err)
   }
}
