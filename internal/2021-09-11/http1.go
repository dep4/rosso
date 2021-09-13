package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "github.com/refraction-networking/utls"
   "net/http"
)

func version(min uint16) []uint16 {
   vs := []uint16{772, 771, 770, 769, 768}
   for k, v := range vs {
      if v == min {
         return vs[:k+1]
      }
   }
   return nil
}

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

const hello1 =
   "771," +
   "49196-49195-49200-49199-159-158-49188-49187-49192-49191-49162-49161-49172-49171-157-156-61-60-53-47-10," +
   "0-10-11-13-35-23-65281,,"

func HelloGolang() *tls.ClientHelloSpec {
   return &tls.ClientHelloSpec{
      CipherSuites:[]uint16{
         0xc02b, 0xc02f, 0xc02c, 0xc030, 0xcca9, 0xcca8, 0xc009, 0xc013, 0xc00a,
         0xc014, 0x9c, 0x9d, 0x2f, 0x35, 0xc012, 0xa, 0x1301, 0x1302, 0x1303,
      },
      CompressionMethods:[]uint8{0x0},
      Extensions:[]tls.TLSExtension{
         &tls.SNIExtension{}, // 0
         &tls.SupportedCurvesExtension{ // 10
            // all fail
            Curves:[]tls.CurveID{0x1d, 0x17, 0x18, 0x19},
         },
         &tls.SupportedPointsExtension{
            // all fail
            SupportedPoints:[]uint8{0x0},
         },
         &tls.SignatureAlgorithmsExtension{ // 13
            SupportedSignatureAlgorithms:[]tls.SignatureScheme{
               0x804, 0x403, 0x807, 0x805, 0x806, 0x401, 0x501, 0x601, 0x503,
               0x603, 0x201, 0x203,
            },
         },
         &tls.SessionTicketExtension{}, // 35
         &tls.UtlsExtendedMasterSecretExtension{}, // 23
         &tls.RenegotiationInfoExtension{Renegotiation:1}, // 65281
      },
      TLSVersMax: 772,
      TLSVersMin: 771,
   }
}
