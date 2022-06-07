# Research

examples/decrypt-cenc: index out of range [0] with length 0

Using this file:

~~~
https://user-images.githubusercontent.com/73562167/172298706-78ff0077-3394-456a-a174-3cd7b7c499e0.mp4
~~~

I get this result:

~~~
> decrypt-cenc -k 22bdb0063805260307ee5045c0f3835a -i enc.mp4 -o dec.mp4
panic: runtime error: index out of range [0] with length 0

goroutine 1 [running]:
main.decryptSamplesInPlace({0xc000150000, 0x5b, 0x883b98?}, {0xc00007dea8, 0x10, 0x20}, 0xc0000202a0)
        D:/Desktop/mp4ff-master/examples/decrypt-cenc/main.go:223 +0x370
main.decryptFragment(0xc0000569c0, {0xc000074930, 0x1, 0x5?}, {0xc00007dea8, 0x10, 0x20})
        D:/Desktop/mp4ff-master/examples/decrypt-cenc/main.go:198 +0x1bf
main.decryptAndWriteSegments({0xc000006070, 0x1, 0x0?}, {0xc000074930, 0x1, 0x2}, {0xc00007dea8, 0x10, 0x20}, {0x883b98, ...})
        D:/Desktop/mp4ff-master/examples/decrypt-cenc/main.go:158 +0x1c9
main.decryptMP4withCenc({0x883b78?, 0xc000006030?}, {0xc00007dea8, 0x10, 0x20}, {0x883b98, 0xc000006038})
        D:/Desktop/mp4ff-master/examples/decrypt-cenc/main.go:146 +0x345
main.start({0x883b78, 0xc000006030}, {0x883b98, 0xc000006038}, {0xc000010200?, 0x85af3d?})
        D:/Desktop/mp4ff-master/examples/decrypt-cenc/main.go:52 +0xd2
main.main()
        D:/Desktop/mp4ff-master/examples/decrypt-cenc/main.go:36 +0x212
~~~

Same input works as expected with other tools:

~~~
mp4decrypt --key 1:22bdb0063805260307ee5045c0f3835a enc.mp4 dec.mp4
~~~
