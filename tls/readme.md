# TLS

## JA3

`ClientHelloSpec` does not store the TLS Version for the initial Client Hello,
only `SupportedVersionsExtension`:

<https://github.com/refraction-networking/utls/blob/0b2885c8/u_common.go#L114-L117>

Then, for the initial Client Hello, TLS version 1.0 is used:

https://github.com/refraction-networking/utls/blob/0b2885c8/conn.go#L944-L948

To solve, we need to make our own struct. Here are some examples:

- https://github.com/dreadl0ck/tlsx/blob/v1.0.0/clientHello.go#L287-L300
- https://github.com/open-ch/ja3/blob/v1.0.1/ja3.go#L15-L24

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
