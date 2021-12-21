# December 20 2021

https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

~~~
CONNECT [2607:f8b0:4000:818::2003]:443 HTTP/1.1
~~~

These work with MITM Proxy:

~~~
curl --cacert mitmproxy-ca-cert.pem -x localhost:8080 https://example.com
curl --cacert mitmproxy-ca.pem -x localhost:8080 https://example.com
~~~

~~~
generate_cert -host 127.0.0.1
generate_cert -host 10.0.2.2
~~~

~~~
.\openssl req -x509 -nodes -days 365 -newkey rsa:4096 -keyout my_site.key `
-out my_site.crt -reqexts v3_req -extensions v3_ca

.\openssl x509 -in my_site.crt -outform der -out my_site.der.crt
~~~
