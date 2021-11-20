package m3u

import (
   "bufio"
   "encoding/json"
   "io"
   "strconv"
   "strings"
)

type ByteRange map[string][]string

func NewByteRange(src io.Reader, prefix string) ByteRange {
   str := make(ByteRange)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if strings.HasPrefix(text, "#") {
         val = text
      } else {
         param := reader{val}
         param.readString(':', '"')
         text = prefix + text
         params, ok := str[text]
         if ok {
            str[text] = append(params, param.str)
         } else {
            str[text] = []string{param.str}
         }
      }
   }
   return str
}

type Stream map[string]string

func Streams(src io.Reader, prefix string) []Stream {
   var dirs []Stream
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      text := buf.Text()
      if strings.HasPrefix(text, "#") {
         dir := newStream(text, prefix)
         if len(dir) >= 1 {
            dirs = append(dirs, dir)
         }
      } else if len(dirs) >= 1 {
         dirs[len(dirs)-1]["URI"] = prefix + text
      }
   }
   return dirs
}

func newStream(src, prefix string) Stream {
   str := reader{src}
   str.readString(':', '"')
   param := make(Stream)
   for {
      key := str.readString('=', '"')
      if key == "" {
         return param
      }
      val := str.readString(',', '"')
      unq, err := strconv.Unquote(val)
      if err == nil {
         val = unq
      }
      if key == "URI" {
         val = prefix + val
      }
      param[key] = val
   }
}

func (s Stream) Struct(val interface{}) error {
   buf, err := json.Marshal(s)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, val)
}

type reader struct {
   str string
}

func (r *reader) readString(sep, enc rune) string {
   out := true
   for k, v := range r.str {
      if v == enc {
         out = !out
      }
      if out && v == sep {
         str := r.str[:k]
         r.str = r.str[k+1:]
         return str
      }
   }
   str := r.str
   r.str = ""
   return str
}
