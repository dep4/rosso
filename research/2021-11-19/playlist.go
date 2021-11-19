package m3u

import (
   "bufio"
   "encoding/csv"
   "io"
   "strconv"
   "strings"
)

// #EXT-X-STREAM-INF
type playlist map[string]map[string]string

func newPlaylist(src io.Reader) (playlist, error) {
   list := make(playlist)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         _, sPar, ok := cutByte(val, ':')
         if ok {
            rCSV := csv.NewReader(strings.NewReader(sPar))
            rCSV.LazyQuotes = true
            aPar, err := rCSV.Read()
            if err != nil {
               return nil, err
            }
            mPar := make(map[string]string)
            for _, par := range aPar {
               key, val, ok := cutByte(par, '=')
               if ok {
                  unq, err := strconv.Unquote(val)
                  if err != nil {
                     mPar[key] = val
                  } else {
                     mPar[key] = unq
                  }
               }
            }
            list[buf.Text()] = mPar
         }
      }
   }
   return list, nil
}
