package m3u

type (
   array []interface{}
   object map[string]interface{}
)

var hls540 = object{
   "#EXT-X-PLAYLIST-TYPE": "VOD",
   "#EXT-X-TARGETDURATION": 5,
   "#EXT-X-VERSION": 4,
   "#EXTINF": array{6, 6, 1},
   "#EXT-X-BYTERANGE": array{
      object{
         "LIMIT OFFSET": "990572@0",
         "URI": "HLS_540.ts",
      },
      object{
         "LIMIT OFFSET": "1023472@990572",
         "URI": "HLS_540.ts",
      },
      object{
         "LIMIT OFFSET": "64108@2014044",
         "URI": "HLS_540.ts",
      },
   },
}
