# TLS

## Modules

- https://godocs.io/github.com/CUCyber/ja3transport
- https://godocs.io/github.com/dreadl0ck/ja3
- https://godocs.io/github.com/open-ch/ja3

## Types

1. https://godocs.io/io#Reader
2. https://godocs.io/github.com/refraction-networking/utls#ClientHelloSpec
3. https://godocs.io/net/http#Transport

## pcap

First install:

https://github.com/emanuele-f/PCAPdroid

Start app, then change from HTTP Server to PCAP File. Then click start, if
prompted to save, choose Downloads. Start Google Chrome and wait for a page to
load. Then stop monitoring, and copy file to computer:

~~~
adb ls /sdcard/Download
adb pull /sdcard/Download/PCAPdroid_22_Oct_15_19_28.pcap
~~~
