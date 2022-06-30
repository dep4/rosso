package os

import (
   "os"
   "path/filepath"
   "strings"
)

var (
   Open = os.Open
   ReadFile = os.ReadFile
   Stdout = os.Stdout
)

func Create(name string) (*os.File, error) {
   var err error
   name, err = clean(name)
   if err != nil {
      return nil, err
   }
   return os.Create(name)
}

func WriteFile(name string, data []byte) error {
   var err error
   name, err = clean(name)
   if err != nil {
      return err
   }
   return os.WriteFile(name, data, os.ModePerm)
}

func clean(name string) (string, error) {
   dir, file := filepath.Split(name)
   if dir != "" {
      err := os.MkdirAll(dir, os.ModePerm)
      if err != nil {
         return "", err
      }
   }
   mapping := func(r rune) rune {
      if strings.ContainsRune(`"*/:<>?\|`, r) {
         return -1
      }
      return r
   }
   file = strings.Map(mapping, file)
   name = filepath.Join(dir, file)
   os.Stderr.WriteString("OpenFile " + name + "\n")
   return name, nil
}
