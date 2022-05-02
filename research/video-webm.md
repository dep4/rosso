# video/webm

## FFmpeg

~~~
.name           = "matroska,webm",
.extensions     = "mkv,mk3d,mka,mks,webm",
.mime_type      = "audio/webm,audio/x-matroska,video/webm,video/x-matroska"
~~~

https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/matroskadec.c

~~~
.name              = "webm",
.mime_type         = "video/webm",
.extensions        = "webm",
~~~

https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/matroskaenc.c

## Mozilla

~~~
.webm
~~~

- <https://developer.mozilla.org/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types>
- https://developer.mozilla.org/Web/Media/Formats/Containers

## MPC-HC

~~~cpp
ADDFMT((_T("webm"),        StrRes(IDS_MFMT_WEBM),        _T("webm")));
~~~

https://github.com/clsid2/mpc-hc/blob/develop/src/mpc-hc/MediaFormats.cpp

## Wikipedia

~~~
.webm
~~~

https://wikipedia.org/wiki/WebM
