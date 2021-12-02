package m3u

import (
   "bufio"
   "io"
   "strconv"
   "strings"
)

type Directive map[string]string

func newDirective(src string) Directive {
   str := reader{src}
   if str.readString(':', '"') == "#EXT-X-I-FRAME-STREAM-INF" {
      return nil
   }
   dir := make(Directive)
   for {
      key := str.readString('=', '"')
      if key == "" {
         return dir
      }
      val := str.readString(',', '"')
      unq, err := strconv.Unquote(val)
      if err == nil {
         val = unq
      }
      dir[key] = val
   }
}

type Playlist map[string]Directive

func NewPlaylist(src io.Reader) Playlist {
   list := make(Playlist)
   var val Directive
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if strings.HasPrefix(text, "#") {
         dir := newDirective(text)
         uri, ok := dir["URI"]
         if ok {
            delete(dir, "URI")
            list[uri] = dir
         } else {
            val = dir
         }
      } else if text != "" {
         list[text] = val
      }
   }
   return list
}
