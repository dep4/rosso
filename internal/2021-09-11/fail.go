package main
import "github.com/refraction-networking/utls"

var fail = &tls.ClientHelloSpec{
   CipherSuites:[]uint16{
      0xc02b, 0xc02f, 0xc02c, 0xc030,
      0xcca9, 0xcca8, 0xc009, 0xc013, 0xc00a, 0xc014, 0x9c, 0x9d, 0x2f, 0x35,
      0xc012, 0xa, 0x1301, 0x1302, 0x1303,
   },
   CompressionMethods:[]uint8{0x0},
   Extensions:[]tls.TLSExtension{
      &tls.SNIExtension{ServerName:""},
      &tls.StatusRequestExtension{},
      &tls.SupportedCurvesExtension{
         Curves:[]tls.CurveID{0x1d, 0x17, 0x18, 0x19},
      },
      &tls.SupportedPointsExtension{
         SupportedPoints:[]uint8{0x0},
      },
      &tls.SignatureAlgorithmsExtension{
         SupportedSignatureAlgorithms:[]tls.SignatureScheme{
            0x804, 0x403, 0x807, 0x805, 0x806, 0x401, 0x501, 0x601, 0x503,
            0x603, 0x201, 0x203,
         },
      },
      &tls.RenegotiationInfoExtension{Renegotiation:1},
      &tls.ALPNExtension{
         []string{"http/1.1"},
      },
      &tls.SCTExtension{},
      &tls.SupportedVersionsExtension{
         Versions:[]uint16{0x304, 0x303, 0x302, 0x301},
      },
      &tls.KeyShareExtension{
         KeyShares:[]tls.KeyShare{
            tls.KeyShare{Group:0x1d},
         },
      },
   },
}
