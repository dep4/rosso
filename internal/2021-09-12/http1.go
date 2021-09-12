package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "github.com/refraction-networking/utls"
   "net/http"
)

func main() {
   for _, site := range sites {
      spec := HelloGolang()
      req, err := http.NewRequest("HEAD", site, nil)
      if err != nil {
         panic(err)
      }
      _, err = ja3.NewTransport(spec).RoundTrip(req)
      fmt.Println(err, site)
   }
}

var sites = []string{
   "https://www.reddit.com",
   "https://github.com",
   "https://nebulance.io",
   "https://stackoverflow.com",
   "https://variety.com",
   "https://vimeo.com",
   "https://www.google.com",
   "https://www.indiewire.com",
   "https://www.wikipedia.org",
   "https://www.youtube.com",
}

func HelloGolang() *tls.ClientHelloSpec {
   return &tls.ClientHelloSpec{
      CipherSuites:[]uint16{
         4865,4866,4867,49195,49199,49196,49200,52393,52392,49171,49172,156,157,
         47,53,
      },
      CompressionMethods:[]uint8{0x0},
      Extensions:[]tls.TLSExtension{
         &tls.SNIExtension{}, // 0
         &tls.UtlsExtendedMasterSecretExtension{}, // 23
         &tls.RenegotiationInfoExtension{Renegotiation:1}, // 65281
         &tls.SupportedCurvesExtension{ // 10
            // all fail
            Curves:[]tls.CurveID{0x1d, 0x17, 0x18, 0x19},
         },
         &tls.SupportedPointsExtension{}, // 11
         &tls.SessionTicketExtension{}, // 35
         &tls.ALPNExtension{ // 16
            AlpnProtocols:[]string{"http/1.1"},
         },
         &tls.StatusRequestExtension{}, // 5
         &tls.SignatureAlgorithmsExtension{ // 13
            SupportedSignatureAlgorithms:[]tls.SignatureScheme{
               0x804, 0x403, 0x807, 0x805, 0x806, 0x401, 0x501, 0x601, 0x503,
               0x603, 0x201, 0x203,
            },
         },
         &tls.SCTExtension{}, // 18
         &tls.KeyShareExtension{ // 51
            KeyShares:[]tls.KeyShare{
               tls.KeyShare{Group:0x1d},
            },
         },
         &tls.PSKKeyExchangeModesExtension{ // 45
            []uint8{tls.PskModeDHE},
         },
         &tls.SupportedVersionsExtension{ // 43
            Versions:[]uint16{772,771},
         },
         &tls.FakeCertCompressionAlgsExtension{}, // 27
         &tls.UtlsPaddingExtension{ // 21
            GetPaddingLen: tls.BoringPaddingStyle,
         },
      },
   }
}

const hello1 =
   "771," +
   "4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53," +
   "0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,,"

func version(min uint16) []uint16 {
   vs := []uint16{772, 771, 770, 769, 768}
   for k, v := range vs {
      if v == min {
         return vs[:k+1]
      }
   }
   return nil
}
