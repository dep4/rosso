# Research

<https://developer.apple.com/documentation/http_live_streaming/http_live_streaming_hls_authoring_specification_for_apple_devices/hls_authoring_specification_for_apple_devices_appendixes>

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

result:

~~~
GROUP-ID="audio-HE-stereo-32_vod-ak-aoc.tv.apple.com"
URI="stream/playlist.m3u8?cc=SI&g=32

GROUP-ID="audio-HE-stereo-64_vod-ak-aoc.tv.apple.com"
URI="stream/playlist.m3u8?cc=SI&g=64

GROUP-ID="audio-stereo-128_vod-ak-aoc.tv.apple.com"
URI="stream/playlist.m3u8?cc=SI&g=128

GROUP-ID="audio-stereo-160_vod-ak-aoc.tv.apple.com"
URI="stream/playlist.m3u8?cc=SI&g=160

GROUP-ID="audio-ac3_vod-ak-aoc.tv.apple.com"
URI="stream/playlist.m3u8?cc=SI&g=384

GROUP-ID="audio-atmos_vod-ak-aoc.tv.apple.com"
URI="stream/playlist.m3u8?cc=SI&g=2448
~~~
