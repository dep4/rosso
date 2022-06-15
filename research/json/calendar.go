package calendar

import (
   "encoding/json"
   "io"
)

type readWriter struct {
   io.Reader
   io.Writer
   n int
}

func (rw *readWriter) Read(p []byte) (int, error) {
   n, err := rw.Reader.Read(p)
   rw.n += n
   return n, err
}

func (rw *readWriter) Write(p []byte) (int, error) {
   n, err := rw.Writer.Write(p)
   rw.n += n
   return n, err
}

type date struct {
   Month int
   Day int
}

func (d *date) ReadFrom(r io.Reader) (int64, error) {
   rw := readWriter{Reader: r}
   err := json.NewDecoder(&rw).Decode(d)
   if err != nil {
      return 0, err
   }
   return int64(rw.n), nil
}

func (d date) WriteTo(w io.Writer) (int64, error) {
   rw := readWriter{Writer: w}
   err := json.NewEncoder(&rw).Encode(d)
   if err != nil {
      return 0, err
   }
   return int64(rw.n), nil
}
