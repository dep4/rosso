# M3U

- https://github.com/oopsguy/m3u8
- https://github.com/ytdl-org/youtube-dl/issues/30491
- https://play.google.com/store/apps/details?id=com.cbs.app
- https://wikipedia.org/wiki/M3U

~~~
yt-dlp --proxy 127.0.0.1:8080 --no-check-certificate `
paramountplus.com/shows/star-trek-prodigy/video/3htV4fvVt4Z8gDZHqlzPOGLSMgcGc_vy/star-trek-prodigy-dreamcatcher
~~~

## #EXT-X-KEY

If IV is missing, then use KEY for both.

https://github.com/oopsguy/m3u8/blob/master/tool/crypt.go#L25-L39

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
