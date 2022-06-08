# Research

- <https://chromium.googlesource.com/chromium/src/+/HEAD/media/cdm/cbcs_decryptor.cc>
- https://github.com/edgeware/mp4ff/issues/150

Workaround:

~~~
packager-win-x64 --enable_raw_key_decryption `
--keys key_id=00000000000000000000000000000000:key=22bdb0063805260307ee5045c0f3835a `
stream=video,in=enc.mp4,output=dec.mp4
~~~
