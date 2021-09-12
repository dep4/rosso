package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "github.com/refraction-networking/utls"
   "net/http"
   "os"
)

func main() {
   req, err := http.NewRequest("GET", "https://www.reddit.com", nil)
   if err != nil {
      panic(err)
   }
   spec := helloGolang()
   res, err := ja3.NewTransport(spec).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   os.Stdout.ReadFrom(res.Body)
   fmt.Printf("%+v\n", res)
}

// Using a function prevents modification.
func helloGolang() *tls.ClientHelloSpec {
   return &tls.ClientHelloSpec{
      CipherSuites:[]uint16{
         0xc02b, 0xc02f, 0xc02c, 0xc030, 0xcca9, 0xcca8, 0xc009, 0xc013, 0xc00a,
         0xc014, 0x9c, 0x9d, 0x2f, 0x35, 0xc012, 0xa, 0x1301, 0x1302, 0x1303,
      },
      CompressionMethods:[]uint8{0x0},
      Extensions:[]tls.TLSExtension{
         &tls.SNIExtension{},
         &tls.StatusRequestExtension{},
         &tls.SupportedCurvesExtension{
            // can omit with reddit? yes.
            Curves:[]tls.CurveID{0x1d, 0x17, 0x18, 0x19},
         },
         &tls.SupportedPointsExtension{
            // can omit with reddit? yes.
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
            AlpnProtocols:[]string{"http/1.1"},
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
}
