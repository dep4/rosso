package m3u

type (
   array []interface{}
   object map[string]interface{}
)

var playlistHLS = object{
   "#EXT-X-VERSION": 4,
   "#EXT-X-STREAM-INF": array{
      object{
         "AVERAGE-BANDWIDTH": 287000,
         "BANDWIDTH": 291000,
         "CLOSED-CAPTIONS": "NONE",
         "CODECS": "avc1.42001e",
         "RESOLUTION": "288x288",
         "URI": "HLS_224.m3u8",
      },
      object{
         "AVERAGE-BANDWIDTH": 497000,
         "BANDWIDTH": 505000,
         "CLOSED-CAPTIONS": "NONE",
         "CODECS": "avc1.42001e",
         "RESOLUTION": "320x320",
         "URI": "HLS_270.m3u8",
      },
      object{
         "AVERAGE-BANDWIDTH": 902000,
         "BANDWIDTH": 920000,
         "CLOSED-CAPTIONS": "NONE",
         "CODECS": "avc1.4d001e",
         "RESOLUTION": "432x432",
         "URI": "HLS_360.m3u8",
      },
      object{
         "AVERAGE-BANDWIDTH": 1352000,
         "BANDWIDTH": 1381000,
         "CLOSED-CAPTIONS": "NONE",
         "CODECS": "avc1.4d001f",
         "RESOLUTION": "640x640",
         "URI": "HLS_540.m3u8",
      },
   },
   "#EXT-X-I-FRAME-STREAM-INF": array{
      object{
         "BANDWIDTH": 102089,
         "CODECS": "avc1.4d001e",
         "RESOLUTION": "432x432",
         "URI": "HLS_360-iframe.m3u8",
      },
      object{
         "BANDWIDTH": 168630,
         "CODECS": "avc1.4d001f",
         "RESOLUTION": "640x640",
         "URI": "HLS_540-iframe.m3u8",
      },
      object{
         "BANDWIDTH": 41018,
         "CODECS": "avc1.42001e",
         "RESOLUTION": "288x288",
         "URI": "HLS_224-iframe.m3u8",
      },
      object{
         "BANDWIDTH": 59248,
         "CODECS": "avc1.42001e",
         "RESOLUTION": "320x320",
         "URI": "HLS_270-iframe.m3u8",
      },
   },
}
