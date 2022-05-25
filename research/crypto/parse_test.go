package crypto

import (
   "fmt"
   "github.com/89z/format/crypto"
   "github.com/refraction-networking/utls"
   "net/http"
   "net/url"
   "strings"
   "testing"
)

const (
   fail = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513-21,29-23-24,0"
   pass = "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"
)

var specPass = []tls.TLSExtension{
   &tls.SNIExtension{ServerName:""},
   &tls.UtlsExtendedMasterSecretExtension{},
   &tls.RenegotiationInfoExtension{Renegotiation:0},
   &tls.SupportedCurvesExtension{Curves:[]tls.CurveID{0x1d, 0x17, 0x18}},
   &tls.SupportedPointsExtension{SupportedPoints:[]uint8{0x0}},
   &tls.GenericExtension{Id:0x23, Data:[]uint8(nil)},
   &tls.ALPNExtension{AlpnProtocols:[]string{"http/1.1"}},
   &tls.StatusRequestExtension{},
   &tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms:[]tls.SignatureScheme{0x403, 0x401}},
   &tls.GenericExtension{Id:0x33, Data:[]uint8(nil)},
   &tls.PSKKeyExchangeModesExtension{Modes:[]uint8{0x1}},
   &tls.SupportedVersionsExtension{Versions:[]uint16{0x303}},
   &tls.GenericExtension{Id:0x15, Data:[]uint8(nil)},
}

var specFail = []tls.TLSExtension{
   &tls.SNIExtension{ServerName:""},
   &tls.UtlsExtendedMasterSecretExtension{},
   &tls.RenegotiationInfoExtension{Renegotiation:0},
   &tls.SupportedCurvesExtension{Curves:[]tls.CurveID{0x1d, 0x17, 0x18}},
   &tls.SupportedPointsExtension{SupportedPoints:[]uint8{0x0}},
   &tls.GenericExtension{Id:0x23, Data:[]uint8(nil)},
   &tls.ALPNExtension{AlpnProtocols:[]string{"http/1.1"}},
   &tls.StatusRequestExtension{},
   &tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms:[]tls.SignatureScheme{0x403, 0x401}},
   &tls.GenericExtension{Id:0x12, Data:[]uint8(nil)},
   &tls.GenericExtension{Id:0x33, Data:[]uint8(nil)},
   &tls.PSKKeyExchangeModesExtension{Modes:[]uint8{0x1}},
   &tls.SupportedVersionsExtension{Versions:[]uint16{0x303}},
   &tls.GenericExtension{Id:0x1b, Data:[]uint8(nil)},
   &tls.GenericExtension{Id:0x4469, Data:[]uint8(nil)},
   &tls.GenericExtension{Id:0x15, Data:[]uint8(nil)},
}

func TestParseJA3(t *testing.T) {
   val := url.Values{
      "Email": {email},
      "Passwd": {password},
      "client_sig": {""},
      "droidguard_results": {""},
   }.Encode()
   specPass, err := crypto.ParseJA3(pass)
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth",
      strings.NewReader(val),
   )
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := crypto.Transport(specPass).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   fmt.Println(res.Status)
}
