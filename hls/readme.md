# HLS

- <https://developer.apple.com/documentation/http_live_streaming/http_live_streaming_hls_authoring_specification_for_apple_devices/hls_authoring_specification_for_apple_devices_appendixes>
- <https://wikipedia.org/wiki/HTTP_Live_Streaming>

## Audio

~~~
audio-HE-stereo-32_vod-ak-aoc.tv.apple.com-Deutsch
aac (HE-AACv2) (mp4a / 0x6134706D), 44100 Hz, stereo, fltp, 30 kb/s

audio-HE-stereo-64_vod-ak-aoc.tv.apple.com-Deutsch
aac (HE-AAC) (mp4a / 0x6134706D), 44100 Hz, stereo, fltp, 61 kb/s

audio-stereo-128_vod-ak-aoc.tv.apple.com-Deutsch
aac (LC) (mp4a / 0x6134706D), 44100 Hz, stereo, fltp, 114 kb/s

audio-stereo-160_vod-ak-aoc.tv.apple.com-Deutsch
aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 142 kb/s

audio-ac3_vod-ak-aoc.tv.apple.com-Deutsch
ac3 (ac-3 / 0x332D6361), 48000 Hz, 5.1(side), fltp, 384 kb/s

audio-atmos_vod-ak-aoc.tv.apple.com-Deutsch
eac3 (ec-3 / 0x332D6365), 48000 Hz, 5.1(side), fltp, 448 kb/s
~~~

## CBC

Why does this:

~~~
#EXT-X-KEY:METHOD=AES-128,URI="https://cbsios-vh.akamaihd.net/i/temp_hd_galle...
~~~

mean CBC?

> An encryption method of AES-128 signals that Media Segments are completely
> encrypted using the Advanced Encryption Standard (AES) [`AES_128`] with a
> 128-bit key, Cipher Block Chaining (CBC)

https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.4

## EXT-X-KEY

If IV is missing, then use KEY for both.
