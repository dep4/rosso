# HLS

- <https://wikipedia.org/wiki/HTTP_Live_Streaming>
- https://godocs.io/net/url#URL.ResolveReference
- https://play.google.com/store/apps/details?id=com.cbs.app

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

https://github.com/oopsguy/m3u8/blob/master/tool/crypt.go#L25-L39
