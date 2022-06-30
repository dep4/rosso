# HLS

## How to get extension?

First look at the master playlist:

~~~
#EXT-X-STREAM-INF:BANDWIDTH=2767539,RESOLUTION=1280x720,CODECS="avc1.4d401f,mp4a.40.2",AUDIO="audio",CLOSED-CAPTIONS="CC"
QualityLevels(2499996)/Manifest(video,format=m3u8-aapl,filter=desktop)
~~~

This does not provide the format, only the codecs. Next look at the segment
playlist:

~~~
#EXTINF:6.037188,no-desc
Fragments(audio_eng_aacl=0,format=m3u8-aapl)
~~~

This does not provide the format, only the codec. Further, the HTTP response
headers cannot be trusted either:

~~~
> curl -I https://s4b3b9a4.ssl.hwcdn.net/files/a8wn4hw/vi/04/07/10427421/hls-mi/s104274210.ts
HTTP/1.1 200 OK
Date: Thu, 30 Jun 2022 18:38:52 GMT
Connection: Keep-Alive
ETag: "1615799593"
Cache-Control: max-age=31536000
Content-Length: 2661328
Content-Type: text/vnd.trolltech.linguist
~~~

which means we must read the response body.
