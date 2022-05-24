# Crypto

- <https://github.com/refraction-networking/utls/blob/master/u_fingerprinter.go>
- https://github.com/golang/net/blob/99673261/http2/hpack/hpack.go#L421-L504
- https://github.com/open-ch/ja3/blob/d0c402d2/parser.go#L277-L335
- https://iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml

## How to get Android JA3?

Check out the `cmd/proxy` folder.

## TLS Handshake Explained - Computerphile

https://youtube.com/watch?v=86cQJ0MMses

## What about Akamai fingerprint?

Paper says HTTP/2 only:

https://blackhat.com/docs/eu-17/materials/eu-17-Shuster-Passive-Fingerprinting-Of-HTTP2-Clients-wp.pdf

Confirmed:

~~~
> curl --http1.1 https://tls.peet.ws/api/clean
{
  "ja3": "771,4866-4867-4865-49196-49200-159-52393-52392-52394-49195-49199-158-49188-49192-107-49187-49191-103-49162-49172-57-49161-49171-51-157-156-61-60-53-47-255,0-11-10-13172-16-22-23-49-13-43-45-51,29-23-30-25-24,0-1-2",
  "ja3_hash": "ba730f97dcd1122e74e65411e68f1b40",
  "akamai": "-",
  "akamai_hash": "-"
}
~~~

I would need to add HTTP/2 support to my existing code:

https://github.com/refraction-networking/utls/blob/9d36ce36/examples/examples.go#L417-L427

So in that case, supporting JA3 is simpler than supporting Akamai.
