# m3u8

## ABC

Both:

~~~
> GET /ausw/slices/933/ab1b22eefb9149afa95e2779c00e9e47/933d34096e1f45dc801d3dd2fef5f598/I00000003.ts?pbs=fe8babf13f4443a4a504bc378549b00d&euid=AC5C4BCC-3F84-4102-9F5B-F639554404F5_000_0_033-05_lf_01-03-00_NA&cloud=aws&cdn=eci&si=0&d=4.096 HTTP/1.1
> Host: x-default-stgec.uplynk.com
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 200 OK
< Accept-Ranges: bytes
< Access-Control-Allow-Origin: *
< Age: 452793
< Cache-Control: no-cache
< cdn-request-id: 96702353280091735075732672073961003582
< Content-Type: application/octet-stream
< Date: Sun, 01 May 2022 23:10:44 GMT
< Etag: "060e14e7f4fc3714a1ab2e7de6e7ca87"
< Expires: Sun, 01 May 2022 23:10:43 GMT
< Last-Modified: Sat, 16 Nov 2019 05:31:55 GMT
< Server: ECAcc (dad/5F8C)
< x-amz-id-2: 2R7LbmQ/l9UxG6VsbM5ww9Bsubs/WZGkuaETfj3MhkR/LQjYhyKr5yhQUDshQHkeZB5VnBP7pRc=
< x-amz-request-id: GWVBD1J486TAZW5V
< x-amz-server-side-encryption: AES256
< X-Cache: HIT
< Content-Length: 1309056
~~~

## NBC

Both:

~~~
> GET /r/NnzsPC/P3SxIJ4UKaJ7,x_U_3FRM6_mL,E8eDNLskrRwg,74sfimEHVQvk,ThnunlAOwYup,EfiqbgG8AVZ8/aHR0cHM6Ly92b2QtbGYtb25lYXBwLXByZC5ha2FtYWl6ZWQubmV0L3Byb2QvdmlkZW8vMWo3L21YZy85MDAwMTk5MzU4L0NFeWlEeW5mQ0d2YUpqZEowdUdVeC8zMDAwa181NDBfaGxzL18xOTI3ODIxNjUwLTNfMDAwMDQudHM?sid=b676b4d6-c127-48cf-8456-93e304e8b272&policy=189081367&date=1651446405877&ip=72.181.23.38&schema=1.0&cid=f1ea39a9-20aa-4a53-a687-1994ac077828&aid=2410887629&dur=3973570&sig=da95c419f72e502a31beffbc1691587f8d2eb904693206133b863b1b6a0fd1cb HTTP/1.1
> Host: redirect.manifest.theplatform.com
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 302 Found
< Date: Sun, 01 May 2022 23:07:16 GMT
< Location: https://vod-lf-oneapp-prd.akamaized.net/prod/video/1j7/mXg/9000199358/CEyiDynfCGvaJjdJ0uGUx/3000k_540_hls/_1927821650-3_00004.ts
< Content-Length: 0
< Server: Jetty(8.1.16.2)

> GET /prod/video/1j7/mXg/9000199358/CEyiDynfCGvaJjdJ0uGUx/3000k_540_hls/_1927821650-3_00004.ts HTTP/1.1
> Host: vod-lf-oneapp-prd.akamaized.net
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 200 OK
< x-amz-id-2: QiHVP0DAfYfZIR5ESGDXkQKiC2CnQvv3R2dL1IEiu8fS4T2e0FAklZf0uaN1fPDh+K1oXy+hb/g=
< x-amz-request-id: NFAXFGZBHS27S6A7
< x-amz-replication-status: COMPLETED
< Last-Modified: Sat, 09 Oct 2021 17:35:48 GMT
< ETag: "3faa651270dd0ae890a56983cdddea1b"
< x-amz-server-side-encryption: AES256
< x-amz-version-id: YqPUejaDWNlYbt5BoXCw4z1IbDGAiylf
< Accept-Ranges: bytes
< Server: AmazonS3
< Content-Length: 2252428
< Cache-Control: max-age=31535971
< Date: Sun, 01 May 2022 23:07:16 GMT
< Connection: keep-alive
< Akamai-Mon-Iucid-Del: 1150039
< Alt-Svc: h3-Q050=":443"; ma=93600,quic=":443"; ma=93600; v="46,43"
< Content-Type: video/MP2T
< Access-Control-Max-Age: 86400
< Access-Control-Allow-Credentials: true
< Access-Control-Expose-Headers: Server,range,hdntl,hdnts,Akamai-Mon-Iucid-Ing,Akamai-Mon-Iucid-Del,Akamai-Request-BC
< Access-Control-Allow-Headers: origin,range,hdntl,hdnts
< Access-Control-Allow-Methods: GET,POST,OPTIONS
< Access-Control-Allow-Origin: *
~~~

## Paramount

## PBS

## CBC
