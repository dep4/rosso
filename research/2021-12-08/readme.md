# December 8 2021

~~~
generate_cert -host 127.0.0.1
generate_cert -host 10.0.2.2
~~~

- <https://wikipedia.org/wiki/Self-signed_certificate>
- https://github.com/jmhodges/howsmyssl/issues/357
- https://medium.com/@j0hnsmith/eavesdrop-on-a-golang-http-client-c4dc49af9d5e
- https://medium.com/@shaneutt/create-sign-x509-certificates-in-golang-8ac4ae49f903
- https://unix.stackexchange.com/questions/208412/how-to-see-list-of-curl

## Done

- HTTP no proxy
- HTTP proxy
- HTTPS no proxy

## To do

HTTPS proxy

~~~
.\openssl req -x509 -nodes -days 365 -newkey rsa:4096 -keyout my_site.key `
-out my_site.crt -reqexts v3_req -extensions v3_ca

.\openssl x509 -in my_site.crt -outform der -out my_site.der.crt
~~~

https://android.stackexchange.com/questions/61540/self-signed-certificate
