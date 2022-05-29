# Dash

## PIFF

https://github.com/asrashley/dashpiff/issues/2

~~~
a2394f525a9b4f14a2446c427c648df4 -language:"Java Properties"

0xa2 0x39 0x4f 0x52 0x5a 0x9b 0x4f 0x14 0xa2 0x44 0x6c 0x42 0x7c 0x64 0x8d 0xf4

"0xa2 0x39 0x4f"
~~~

https://github.com/truedread/pymp4decrypt/issues/1

## amcplus.com/shows/orphan-black/episodes/season-1-instinct--1011152

This one is encrypted. Key:

~~~
a66a5603545ad206c1a78e160a6710b1
~~~

Get init and first segment:

~~~
amc -b 1011152 -f 1 -g 0
~~~

Try to decrypt:

~~~
.\decrypt-cenc `
-k a66a5603545ad206c1a78e160a6710b1 `
-i 1011152-enc.mp4 `
-o 1011152-dec.mp4
~~~

Then sanity check:

~~~
mp4decrypt `
--key 1:a66a5603545ad206c1a78e160a6710b1 `
1011152-enc.mp4 `
1011152-dec.mp4
~~~

FFmpeg works too:

~~~
ffmpeg `
-decryption_key a66a5603545ad206c1a78e160a6710b1 `
-i 1011152-enc.mp4 `
-c copy `
1011152-dec.mp4
~~~

## bbcamerica.com/shows/killing-eve/episodes/season-4-just-dunk-me--1052529

~~~
amc -b 1052529 -f 1
~~~

result is not encrypted.
