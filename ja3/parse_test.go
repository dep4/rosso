package ja3

import (
   "net/http"
   "testing"
)

const test =
   "771," +
   "4866-4867-4865-49196-49200-49195," +
   "0-11-10-16-22-23-49-13-43-45-51-21," +
   "29-23-1035-25-24," +
   "0-1-2"

func TestParse(t *testing.T) {
   spec, err := Parse(test)
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest("HEAD", "https://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   if _, err := NewTransport(spec).RoundTrip(req); err != nil {
      t.Fatal(err)
   }
}
