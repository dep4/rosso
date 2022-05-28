# Dash

Using this:

https://bbcamerica.com/shows/killing-eve/episodes/season-4-just-dunk-me--1052529

~~~
680a46ebd6cf2b9a6a0b05a24dcf944a
~~~

then:

~~~
http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/b7173ce1-3126-4042-a0d1-f8cf8dacd94e/xd9/default_audio128_5_en_main/init0.m4f
http://redirector.playback.us-east-1.prod.deploys.brightcove.com/v1/6245817279001/b7173ce1-3126-4042-a0d1-f8cf8dacd94e/xd9/default_audio128_5_en_main/segment0.m4f
~~~

then:

~~~
.\decrypt-cenc.exe `
-k 680a46ebd6cf2b9a6a0b05a24dcf944a `
-i 1052529-enc.mp4 `
-o 1052529-dec.mp4
~~~
