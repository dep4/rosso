# Research

For MPC-HC, the video and audio files have to have the same name. This means we
need to use extension for both files, and the extension need to be different.
CBC uses MPEG-TS for both files:

~~~
Input #0, mpegts, from 'CBC_DOWNTON_ABBEY_S01E05-v':
  Stream #0:0[0x12c]: Video: h264 (Constrained Baseline) ([27][0][0][0] /

Input #0, mpegts, from 'CBC_DOWNTON_ABBEY_S01E05-a':
  Stream #0:0[0x12d]: Audio: aac (LC) ([15][0][0][0] / 0x000F), 44100 Hz,
~~~

Ideally we would use `.tsv` and `.tsa`, but these are not supported by FFmpeg:

~~~
[NULL @ 0000023376e1bb00] Unable to find a suitable output format for 'out.tsv'
[NULL @ 000002630ebfbb00] Unable to find a suitable output format for 'out.tsa'
~~~

Or Mozilla:

<https://developer.mozilla.org/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types>

Or MPC-HC:

https://github.com/clsid2/mpc-hc/issues/1697

I could work around or ignore these issues, but I think the better option is to
look at the codecs instead. Starting with video:

<https://wikipedia.org/wiki/Advanced_Video_Coding>

~~~
3g2
3ga
3gp
3gp2
3gpp
aac
ac3
aif
aifc
aiff
alac
amr
amv
aob
ape
apl
asf
asx
au
avi
bdmv
bik
caf
cda
cue
dav
divx
dsa
dsf
dsm
dss
dsv
dts
dtshd
dtsma
dv
evo
f4v
flac
flc
fli
flic
flv
hdmov
ifo
ivf
m1a
m1v
m2a
m2p
m2t
m2ts
m2v
m3u
m3u8
m4a
m4b
m4r
m4v
mid
midi
mk3d
mka
mkv
mlp
mov
mp2
mp2v
mp3
mp4
mp4v
mpa
mpc
mpcpl
mpe
mpeg
mpg
mpl
mpls
mpv2
mpv4
mts
mxf
ofr
ofs
oga
ogg
ogm
ogv
opus
pls
pva
ra
ram
rar
rec
rm
rmi
rmm
rmvb
rp
rpm
rt
smi
smil
smk
snd
spx
ssif
swf
tak
thd
tp
trp
ts
tta
vob
wav
wax
weba
webm
wm
wma
wmp
wmv
wmx
wv
wvx
~~~

https://github.com/clsid2/mpc-hc/blob/develop/src/mpc-hc/MediaFormats.cpp
