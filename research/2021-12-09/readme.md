# December 9 2021

These work with MITM Proxy:

~~~
curl --cacert mitmproxy-ca-cert.pem -x localhost:8080 https://example.com
curl --cacert mitmproxy-ca.pem -x localhost:8080 https://example.com
~~~

Search:

~~~
mitm language:go size:0..10 stars:>0
~~~

- https://github.com/Bren2010/mitm/issues/1
- https://github.com/petethepig/mitm/issues/2
- https://github.com/vinhjaxt/vitm-proxy/issues/1

These work:

- https://github.com/JetDoughnut/ch
- https://github.com/agatan/curlize
