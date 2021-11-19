package m3u

import (
   "bufio"
   "io"
   "strings"
)

var hlsPlaylist = stream{
   "HLS_540.m3u8": {
      "AVERAGE-BANDWIDTH": "1352000",
      "BANDWIDTH": "1381000",
      "CLOSED-CAPTIONS": "NONE",
      "CODECS": "avc1.4d001f",
      "RESOLUTION": "640x640",
   },
}

var hls540 = byteRange{
   "HLS_540.ts": {
      "1023472@990572",
      "64108@2014044",
      "990572@0",
   },
}

// #EXT-X-STREAM-INF
type stream map[string]map[string]string

func newStream() stream {
   return make(stream)
}

// #EXT-X-BYTERANGE
type byteRange map[string][]string

func newByteRange() byteRange {
   return make(byteRange)
}

////////////////////////////////////////////////////////////////////////////////

func Parse(src io.Reader, prefix string) map[string][]string {
   list := make(map[string][]string)
   var val string
   buf := bufio.NewScanner(src)
   for buf.Scan() {
      if strings.HasPrefix(buf.Text(), "#") {
         val = buf.Text()
      } else {
         _, param, ok := cutByte(val, ':')
         if ok {
            text := prefix + buf.Text()
            params, ok := list[text]
            if ok {
               list[text] = append(params, param)
            } else {
               list[text] = []string{param}
            }
         }
      }
   }
   return list
}

func cutByte(s string, sep byte) (string, string, bool) {
   i := strings.IndexByte(s, sep)
   if i == -1 {
      return s, "", false
   }
   return s[:i], s[i+1:], true
}
