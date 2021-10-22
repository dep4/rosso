# October 21 2021

~~~
netsh trace start capture=yes tracefile=NetTrace.etl
netsh trace stop
~~~

Or:

~~~
pktmon start --etw -p 0
pktmon stop
~~~

- https://blog.rmilne.ca/2016/08/11/network-monitor-filter-examples
- https://github.com/microsoft/etl2pcapng/issues/48
- https://techcommunity.microsoft.com/t5/iis-support-blog/capture/ba-p/376503
