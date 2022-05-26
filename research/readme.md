# Research

~~~
pass
mp4decrypt --key 1:680a46ebd6cf2b9a6a0b05a24dcf944a enc.mp4 dec.mp4

fail
ffmpeg -decryption_key 680a46ebd6cf2b9a6a0b05a24dcf944a `
-i enc.mp4 -c copy dec.mp4

fail
ffmpeg -cenc_decryption_key 680a46ebd6cf2b9a6a0b05a24dcf944a `
-i enc.mp4 -c copy dec.mp4
~~~

- https://github.com/GyanD/codexffmpeg/issues/53
- https://github.com/edgeware/mp4ff/blob/master/examples/decrypt-cenc/main.go
- https://github.com/edgeware/mp4ff/issues/146
