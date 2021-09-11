# Internal

~~~
&tls.SNIExtension{ServerName:""}
&tls.SupportedPointsExtension{SupportedPoints:[]uint8{0x0, 0x1, 0x2}}
&tls.SupportedCurvesExtension{Curves:[]tls.CurveID{0x1d, 0x17, 0x1e, 0x19, 0x18}}
&tls.NPNExtension{NextProtos:[]string(nil)}
&tls.ALPNExtension{AlpnProtocols:[]string{"h2", "http/1.1"}}
&tls.GenericExtension{Id:0x16, Data:[]uint8{}}
&tls.UtlsExtendedMasterSecretExtension{}
&tls.GenericExtension{Id:0x31, Data:[]uint8{}}
&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms:[]tls.SignatureScheme{0x403, 0x503, 0x603, 0x807, 0x808, 0x809, 0x80a, 0x80b, 0x804, 0x805, 0x806, 0x401, 0x501, 0x601, 0x303, 0x203, 0x301, 0x201, 0x302, 0x202, 0x402, 0x502, 0x602}}
&tls.SupportedVersionsExtension{Versions:[]uint16{0x304, 0x303, 0x302, 0x301}}
&tls.PSKKeyExchangeModesExtension{Modes:[]uint8{0x1}}
&tls.KeyShareExtension{KeyShares:[]tls.KeyShare{tls.KeyShare{Group:0x1d, Data:[]uint8(nil)}}}
&tls.UtlsPaddingExtension{PaddingLen:0, WillPad:false, GetPaddingLen:(func(int) (int, bool))(0x2659e0)}
~~~
