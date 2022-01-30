package m3u

import (
   "bufio"
   "io"
   "strings"
)

type format struct {
   resolution, codecs, uri string
}

func decode(src io.Reader) []format {
   var forms []format
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      pair := strings.SplitN(buf.Text(), ":", 2)
      if len(pair) == 2 && pair[0] == "#EXT-X-STREAM-INF" {
         var form format
         for _, property := range strings.Split(pair[1], ",") {
            pair := strings.SplitN(property, "=", 2)
            if len(pair) == 2 && pair[0] == "RESOLUTION" {
               form.resolution = pair[1]
            }
         }
         buf.Scan()
         form.uri = buf.Text()
         forms = append(forms, form)
      }
   }
   return forms
}
