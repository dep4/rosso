package m3u

import (
   "bufio"
   "encoding/csv"
   "io"
   "strconv"
   "strings"
)

func cutByte(s string, sep byte) (string, string, bool) {
   i := strings.IndexByte(s, sep)
   if i == -1 {
      return s, "", false
   }
   return s[:i], s[i+1:], true
}

// #EXT-X-STREAM-INF
type Playlist map[string]map[string]string

func NewPlaylist(src io.Reader, prefix string) (Playlist, error) {
   list := make(Playlist)
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
            list[prefix + buf.Text()] = mPar
         }
      }
   }
   return list, nil
}

// #EXT-X-BYTERANGE
type Stream map[string][]string

func NewStream(src io.Reader, prefix string) Stream {
   str := make(Stream)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         _, param, ok := cutByte(val, ':')
         if ok {
            text := prefix + buf.Text()
            params, ok := str[text]
            if ok {
               str[text] = append(params, param)
            } else {
               str[text] = []string{param}
            }
         }
      }
   }
   return str
}
