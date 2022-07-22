package os

import (
   "os"
   "path/filepath"
   "strings"
)

func Create(name string) (*os.File, error) {
   os.Stderr.WriteString("Create " + name + "\n")
   return os.Create(name)
}

func Rename(old_path, new_path string) error {
   os.Stderr.WriteString("Rename " + new_path + "\n")
   return os.Rename(old_path, new_path)
}

func WriteFile(name string, data []byte) error {
   os.Stderr.WriteString("WriteFile " + name + "\n")
   return os.WriteFile(name, data, os.ModePerm)
}

type Cleaner struct {
   name string
}

func Clean(dir, file string) Cleaner {
   mapping := func(r rune) rune {
      if strings.ContainsRune(`"*/:<>?\|`, r) {
         return -1
      }
      return r
   }
   var c Cleaner
   c.name = strings.Map(mapping, file)
   c.name = filepath.Join(dir, c.name)
   return c
}

func (c Cleaner) Create() (*os.File, error) {
   os.Stderr.WriteString("Create " + c.name + "\n")
   return os.Create(c.name)
}
