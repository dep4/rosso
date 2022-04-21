package json

import (
   "fmt"
   "os"
   "testing"
)

func TestFacebook(t *testing.T) {
   file, err := os.Open("ignore/facebook.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"\u0040context"`)
   scan.Scan()
   var object struct {
      DateCreated string
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}

func TestPbsFrontline(t *testing.T) {
   file, err := os.Open("ignore/pbs-frontline.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"@context"`)
   scan.Scan()
   var object struct {
      Graph []struct {
         EmbedURL string
      } `json:"@graph"`
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}

func TestPbsMasterpiece(t *testing.T) {
   file, err := os.Open("ignore/pbs-masterpiece.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`"https://video.`)
   scan.Scan()
   var object string
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Println(object)
}

func TestPbsNature(t *testing.T) {
   file, err := os.Open("ignore/pbs-nature.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"preview"`)
   scan.Scan()
   var object struct {
      Full_Length map[string]struct {
         Video_Iframe string
      }
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}

func TestPbsNova(t *testing.T) {
   file, err := os.Open("ignore/pbs-nova.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"props"`)
   scan.Scan()
   var object struct {
      Query struct {
         Video string
      }
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}

func TestPbsVideo(t *testing.T) {
   file, err := os.Open("ignore/pbs-video.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte("{\n  \"@context\"")
   scan.Scan()
   var object struct {
      Video struct {
         ContentURL string
      }
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}

func TestPbsWidget(t *testing.T) {
   file, err := os.Open("ignore/pbs-widget.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"availability"`)
   scan.Scan()
   var object struct {
      Encodings []string
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}
