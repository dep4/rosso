package m3u

import (
   "path"
   "strconv"
   "text/scanner"
)

func Masters(src string, fn openFile) ([]Master, error) {
   file, err := fn(src)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   var (
      buf scanner.Scanner
      mass []Master
   )
   buf.Init(file)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas Master
         for buf.Scan() != '\n' {
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               val, err := strconv.ParseInt(buf.TokenText(), 10, 64)
               if err != nil {
                  return nil, err
               }
               mas.Bandwidth = val
            }
         }
         scanLines(&buf)
         buf.Scan()
         mas.URI = buf.TokenText()
         if !isAbsPath(mas.URI) {
            // FIXME fails with Windows
            mas.URI = path.Dir(src) + buf.TokenText()
         }
         mass = append(mass, mas)
      }
   }
   return mass, nil
}
