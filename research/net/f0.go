package net

import (
   "io"
   "net/url"
)

type Values struct {
   url.Values
}

type writer struct {
   n int
   w io.Writer
}

func (w *writer) Write(p []byte) (int, error) {
   n, err := w.w.Write(p)
   w.n += n
   return n, err
}
