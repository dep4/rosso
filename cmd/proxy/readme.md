# Proxy

To capture Android TLS handshake, go to Chrome App info, then Storage, then
MANAGE SPACE, then CLEAR ALL DATA, then OK. Then start the server, and go to
Android Emulator Extended Controls. Choose Manual proxy configuration, then
enter:

~~~
127.0.0.1:8080
~~~

and click Apply. Then start Android Chrome.

- https://android.stackexchange.com/questions/243184/capture-tls-handshake
- https://github.com/spritesprite/proxychannel
- https://unix.stackexchange.com/questions/208412/how-to-see-list-of-curl
