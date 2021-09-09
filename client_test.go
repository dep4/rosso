package ja3
import "testing"

var tests = []string{
   "769,47-53-5-10-49161-49162-49171-49172-50-56-19-4,0-10-11,23-24-25,0",
   "769,47,0-10-11-22,23-24-25,0",
}

func TestClient(t *testing.T) {
   for _, test := range tests {
      _, err := StringToSpec(test)
      if err != nil {
         t.Fatal(err)
      }
   }
}
