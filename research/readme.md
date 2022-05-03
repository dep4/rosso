# Research

For MPC-HC, the video and audio files have to have the same name. This means we
need to use extension for both files, and the extension need to be different.
CBC uses MPEG-TS for both files:

~~~
Input #0, mpegts, from 'CBC_DOWNTON_ABBEY_S01E05-v':
  Stream #0:0[0x12c]: Video: h264 (Constrained Baseline) ([27][0][0][0] /

Input #0, mpegts, from 'CBC_DOWNTON_ABBEY_S01E05-a':
  Stream #0:0[0x12d]: Audio: aac (LC) ([15][0][0][0] / 0x000F), 44100 Hz,
~~~

Ideally we would use `.tsv` and `.tsa`, but these are not supported by FFmpeg:

~~~
[NULL @ 0000023376e1bb00] Unable to find a suitable output format for 'out.tsv'
[NULL @ 000002630ebfbb00] Unable to find a suitable output format for 'out.tsa'
~~~

Or Mozilla:

<https://developer.mozilla.org/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types>

Or MPC-HC:

https://github.com/clsid2/mpc-hc/issues/1697

If we look at PBS:

~~~
GET /videos/nature/1803a543-5d57-41f6-81bf-d199448d45f6/2000281849/hd-16x9-mezzanine-1080p/naat4008_r-hls-16x9-1080p-234p-145k_00001.ts HTTP/1.1
Host: ga.video.cdn.pbs.org

GET /videos/nature/1803a543-5d57-41f6-81bf-d199448d45f6/2000281849/hd-16x9-mezzanine-1080p/naat4008_r-hls-16x9-1080pAudio%20Selector%201_00001.aac HTTP/1.1
Host: ga.video.cdn.pbs.org
~~~

they are using `.ts` for video and `.aac` for audio. FFmpeg can transcode `.ts`
to `.aac`:

~~~
Input #0, mpegts, from 'CBC_DOWNTON_ABBEY_S01E05':
  Duration: 00:49:39.93, start: 0.000000, bitrate: 200 kb/s
  Program 1
  Stream #0:0[0x12d]: Audio: aac (LC) ([15][0][0][0] / 0x000F), 44100 Hz, stereo, fltp, 192 kb/s
Output #0, adts, to '.aac':
  Metadata:
    encoder         : Lavf59.16.100
  Stream #0:0: Audio: aac (LC) ([15][0][0][0] / 0x000F), 44100 Hz, stereo, fltp, 192 kb/s
Stream mapping:
  Stream #0:0 -> #0:0 (copy)
~~~

`.aac` is documented by Mozilla:

<https://developer.mozilla.org/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types>

`.aac` is supported by MPC-HC:

https://github.com/clsid2/mpc-hc/blob/develop/src/mpc-hc/MediaFormats.cpp

`.aac` is documented by Wikipedia:

<https://wikipedia.org/wiki/Advanced_Audio_Coding>
