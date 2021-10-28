# October 27 2021

Go defaults:

~~~
&tls.SNIExtension{ServerName:""}
&tls.StatusRequestExtension{}
&tls.SupportedCurvesExtension{Curves:[]tls.CurveID{0x1d, 0x17, 0x18, 0x19}}
&tls.SupportedPointsExtension{SupportedPoints:[]uint8{0x0}}
&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms:[]tls.SignatureScheme{0x804, 0x403, 0x807, 0x805, 0x806, 0x401, 0x501, 0x601, 0x503, 0x603, 0x201, 0x203}}
&tls.RenegotiationInfoExtension{Renegotiation:1}
&tls.ALPNExtension{AlpnProtocols:[]string{"h2", "http/1.1"}}
&tls.SCTExtension{}
&tls.SupportedVersionsExtension{Versions:[]uint16{0x304, 0x303, 0x302, 0x301}}
&tls.KeyShareExtension{KeyShares:[]tls.KeyShare{tls.KeyShare{Group:0x1d, Data:[]uint8(nil)}}}
~~~
