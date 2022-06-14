package main

import (
   "io"
   "os"
   "strings"
)

type BlockMode struct {
   io.Reader
   protected []byte
   clear []byte
}

func (b *BlockMode) Read(p []byte) (int, error) {
   protected, err := b.Reader.Read(p)
   // add to protected
   b.protected = append(b.protected, p[:protected]...)
   cut := len(b.protected) - len(b.protected) % 3
   // add to clear
   b.clear = append(b.clear, b.protected[:cut]...)
   // remove from protected
   b.protected = b.protected[cut:]
   // add to out
   clear := copy(p, b.clear)
   // remove from clear
   b.clear = b.clear[clear:]
   // return
   return clear, err
}

func main() {
   var mode BlockMode
   mode.Reader = strings.NewReader("0123456789")
   os.Stdout.ReadFrom(&mode)
}
