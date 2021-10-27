# TLS

First install:

https://github.com/emanuele-f/PCAPdroid

Start app, then change from HTTP Server to PCAP File. Then click start, if
prompted to save, choose Downloads. Start Google Chrome and wait for a page to
load. Then stop monitoring, and copy file to computer:

~~~
adb ls /sdcard/Download
adb pull /sdcard/Download/PCAPdroid_22_Oct_15_19_28.pcap
~~~

- https://github.com/open-ch/ja3/blob/d0c402d2/parser.go#L277-L335
- https://iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml
