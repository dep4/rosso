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

## CBC

Video:

~~~
> GET /0f73fb9d-87f0-4577-81d1-e6e970b89a69/CBC_DOWNTON_ABBEY_S01E05.ism/QualityLevels(2500080)/Fragments(video=0,format=m3u8-aapl) HTTP/1.1
> Host: cbcrcott-gem.akamaized.net
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 200 OK
< Pragma: IISMS/6.0,IIS Media Services Premium by Microsoft
< Content-Type: video/mp2t
< ETag: "0x8D8CA2267051B64"
< Server: Microsoft-IIS/10.0 IISMS/6.0
< x-ms-streaming-duration: video=6006
< X-Content-Type-Options: nosniff
< Content-Length: 1946560
< Cache-Control: private, max-age=17731814
< Expires: Wed, 23 Nov 2022 07:03:01 GMT
< Date: Mon, 02 May 2022 01:32:47 GMT
< Connection: keep-alive
< Akamai-Mon-Iucid-Del: 550858
< Alt-Svc: h3-Q050=":443"; ma=93600,quic=":443"; ma=93600; v="46,43"
< Set-Cookie: akaalb_LB-SrcAv-Toutv=~op=SRC_RcAvToutv_lb:Standard-East|~rv=90~m=Standard-East:0|~os=549dc91727a25e1d5313306552a5fc3a~id=fa12fbb303d1e2fea10ef6e432f9e74e; path=/; HttpOnly; Secure; SameSite=None
< Access-Control-Max-Age: 86400
< Access-Control-Allow-Credentials: true
< Access-Control-Expose-Headers: Server,range,hdntl,hdnts,Akamai-Mon-Iucid-Ing,Akamai-Mon-Iucid-Del,Akamai-Request-BC
< Access-Control-Allow-Headers: origin,range,hdntl,hdnts
< Access-Control-Allow-Methods: GET,POST,OPTIONS
< Access-Control-Allow-Origin: *
~~~

Audio:

~~~
> GET /0f73fb9d-87f0-4577-81d1-e6e970b89a69/CBC_DOWNTON_ABBEY_S01E05.ism/QualityLevels(192000)/Fragments(audio_eng_aacl=0,format=m3u8-aapl) HTTP/1.1
> Host: cbcrcott-gem.akamaized.net
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 200 OK
< Pragma: IISMS/6.0,IIS Media Services Premium by Microsoft
< Content-Type: video/mp2t
< ETag: "0x8D8CA233CC10B47"
< Server: Microsoft-IIS/10.0 IISMS/6.0
< x-ms-streaming-duration: audio=6037
< X-Content-Type-Options: nosniff
< Content-Length: 153792
< X-EdgeConnect-MidMile-RTT: 1
< X-EdgeConnect-Origin-MEX-Latency: 311
< X-EdgeConnect-MidMile-RTT: 21
< X-EdgeConnect-Origin-MEX-Latency: 311
< Cache-Control: private, max-age=17728367
< Expires: Wed, 23 Nov 2022 06:06:38 GMT
< Date: Mon, 02 May 2022 01:33:51 GMT
< Connection: keep-alive
< Akamai-Mon-Iucid-Del: 550858
< Alt-Svc: h3-Q050=":443"; ma=93600,quic=":443"; ma=93600; v="46,43"
< Set-Cookie: akaalb_LB-SrcAv-Toutv=~op=SRC_RcAvToutv_lb:Standard-East|~rv=35~m=Standard-East:0|~os=549dc91727a25e1d5313306552a5fc3a~id=de444bf7c9f74037fc47085286267fbb; path=/; HttpOnly; Secure; SameSite=None
< Access-Control-Max-Age: 86400
< Access-Control-Allow-Credentials: true
< Access-Control-Expose-Headers: Server,range,hdntl,hdnts,Akamai-Mon-Iucid-Ing,Akamai-Mon-Iucid-Del,Akamai-Request-BC
< Access-Control-Allow-Headers: origin,range,hdntl,hdnts
< Access-Control-Allow-Methods: GET,POST,OPTIONS
< Access-Control-Allow-Origin: *
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

Both:

~~~
> GET /i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/segment6_3_av.ts?null=0&id=AgBItRcmFy81SGkeb2Lj9yw5Dv67QK6XzSdkg4bVbdovy8FQlOh1cjMYr%2f6IPgZORkdzy1dt3vTjUQ%3d%3d&hdntl=exp=1651535849~acl=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_*~data=hdntl~hmac=7627d6e8f3fd297066625128a81df1018c5d4917aa99d56cac4dd2189a06acfd HTTP/1.1
> Host: cbsios-vh.akamaihd.net
> User-Agent: Go-http-client/1.1
> Accept: */*

< HTTP/1.1 200 OK
< Server: AkamaiGHost
< Mime-Version: 1.0
< Content-Type: video/MP2T
< Content-Length: 2772448
< Pragma: no-cache
< Cache-Control: no-store
< Expires: Sun, 01 May 2022 23:58:50 GMT
< Date: Sun, 01 May 2022 23:58:50 GMT
< Connection: keep-alive
< x-cdn: Akamai
< Access-Control-Allow-Headers: *
< Access-Control-Expose-Headers: *
< Access-Control-Allow-Methods: GET, HEAD, OPTIONS
< Access-Control-Allow-Origin: *
< Set-Cookie: _alid_=1xxteWrwr7qqOl2ZsB/kLA==; path=/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2012/09/12/41581439/CBS_MELROSE_PLACE_001_SD_prores_78930_,503,4628,3128,2228,1628,848,000.mp4.csmil/; domain=cbsios-vh.akamaihd.net
~~~

## PBS

Video:

~~~
> GET /videos/nature/1803a543-5d57-41f6-81bf-d199448d45f6/2000281849/hd-16x9-mezzanine-1080p/naat4008_r-hls-16x9-1080p-234p-145k_00001.ts HTTP/1.1
> Host: ga.video.cdn.pbs.org
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 200 OK
< Content-Type: video/MP2T
< Content-Length: 118816
< Connection: keep-alive
< Date: Thu, 24 Feb 2022 00:31:46 GMT
< Last-Modified: Tue, 22 Feb 2022 15:01:24 GMT
< ETag: "60cacbdae4b6b8efcf533356b06cd3aa"
< x-amz-version-id: 2tJMNqe_E1Jfz6o_2OFGjLUNYy5R6DWa
< Accept-Ranges: bytes
< Server: AmazonS3
< X-Cache: Hit from cloudfront
< Via: 1.1 0b411dbb186753d7d6bc75c4c3de15a0.cloudfront.net (CloudFront)
< X-Amz-Cf-Pop: DFW3-C1
< X-Amz-Cf-Id: OAJAi6NvpcchSoJ7_YHLz9tue08CHD5p1TBYOPAFKOvZy2Za1d7sow==
< Age: 5792266
~~~

Audio:

~~~
> GET /videos/nature/1803a543-5d57-41f6-81bf-d199448d45f6/2000281849/hd-16x9-mezzanine-1080p/naat4008_r-hls-16x9-1080pAudio%20Selector%201_00001.aac HTTP/1.1
> Host: ga.video.cdn.pbs.org
> User-Agent: curl/7.78.0
> Accept: */*

< HTTP/1.1 200 OK
< Content-Type: audio/aac
< Content-Length: 121512
< Connection: keep-alive
< Date: Thu, 24 Feb 2022 00:30:32 GMT
< Last-Modified: Tue, 22 Feb 2022 15:01:24 GMT
< ETag: "70347b7861d0610593dc6fd1777006e5"
< x-amz-version-id: z8BBj76loLAa7TAomBjwSEIeV.gm4U4z
< Accept-Ranges: bytes
< Server: AmazonS3
< X-Cache: Hit from cloudfront
< Via: 1.1 ac72c6dfa21a23a34396bec16bd466a6.cloudfront.net (CloudFront)
< X-Amz-Cf-Pop: DFW56-P5
< X-Amz-Cf-Id: rluS_tiQkLQAj7SBgqZ_DnzhCPzIFW4s1vWwHMO0jQPkDpI4E5GCyA==
< Age: 5792227
~~~
