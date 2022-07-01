package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "mpd/amc.mpd",
   "mpd/paramount-lang.mpd",
   "mpd/paramount-role.mpd",
   "mpd/roku.mpd",
}

func Test_Audio(t *testing.T) {
   for i, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(test)
      for _, rep := range med.Representations().Audio() {
         fmt.Println(rep)
      }
   }
}

func Test_Video(t *testing.T) {
   for i, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(test)
      for _, rep := range med.Representations().Video() {
         fmt.Println(rep)
      }
   }
}

func Test_Representations(t *testing.T) {
   for i, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(test)
      for _, rep := range med.Representations() {
         fmt.Println(rep)
      }
   }
}

func Test_Ext(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      fmt.Println(test)
      for _, rep := range med.Representations() {
         fmt.Printf("%q\n", rep.Ext())
      }
   }
}
