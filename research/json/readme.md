Implement io.WriterTo with JSON

I found this cool interface recently, `io.WriterTo`:

https://godocs.io/io#WriterTo

I would like to implement it for some JSON objects. I was able to make this:

~~~
package calendar

import (
   "bytes"
   "encoding/json"
   "io"
)

type date struct {
   Month int
   Day int
}

func (d date) WriteTo(w io.Writer) (int64, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(d)
   if err != nil {
      return 0, err
   }
   return buf.WriteTo(w)
}
~~~

but I think its not ideal, as it makes a copy of the object in memory, before
sending to the Writer. Is it possible to write directly, but also know how many
bytes were written?
