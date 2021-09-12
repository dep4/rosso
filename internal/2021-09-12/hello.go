package main

import (
   "fmt"
   "github.com/89z/parse/ja3"
   "github.com/refraction-networking/utls"
   "net/http"
)

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

/*
"771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,,"
"771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
"771,4865-4866-4867-49196-49195-52393-49200-49199-52392-49188-49187-49162-49161-49192-49191-49172-49171-157-156-61-60-53-47-49160-49170-10,0-23-65281-10-11-16-5-13-18-51-45-43-21,,"
"771,4865-4866-49196-49195-49200-157-49198-49202-159-163-49199-156-49197-49201-158-162-49188-49192-61-49190-49194-107-106-49162-49172-53-49157-49167-57-56-49187-49191-60-49189-49193-103-64-49161-49171-47-49156-49166-51-50-255,0-5-10-11-13-50-23-43-45-51,,"
"771,4866-4867-4865-49199-49195-49200-49196-158-49191-103-49192-107-163-159-52393-52392-52394-49327-49325-49315-49311-49245-49249-49239-49235-162-49326-49324-49314-49310-49244-49248-49238-49234-49188-106-49187-64-49162-49172-57-56-49161-49171-51-50-157-49313-49309-49233-156-49312-49308-49232-61-60-53-47-255,0-11-10-35-22-23-13-43-45-51,,"
"771,49172-53-49171-47-49196-49195-49200-157-49199-156,10-11-13-23-0-65281,,"
"771,49196-49195-49200-49199-159-158-49188-49187-49192-49191-49162-49161-49172-49171-157-156-61-60-53-47-10,0-10-11-13-35-23-65281,,"
"771,52393-52392-49195-49199-49196-49200-49171-49172-156-157-47-53-10,65281-0-23-35-13-5-18-16-11-10-27,29-23-24,0"
*/

func HelloGolang() *tls.ClientHelloSpec {
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
            AlpnProtocols:[]string{"http/1.1"},
         },
         &tls.SCTExtension{},
         &tls.SupportedVersionsExtension{
            Versions: version(771),
         },
         &tls.KeyShareExtension{
            KeyShares:[]tls.KeyShare{
               tls.KeyShare{Group:0x1d},
            },
         },
      },
   }
}

func version(min uint16) []uint16 {
   vs := []uint16{772, 771, 770, 769, 768}
   for k, v := range vs {
      if v == min {
         return vs[:k+1]
      }
   }
   return nil
}
