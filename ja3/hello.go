package ja3
import "github.com/refraction-networking/utls"

func ClientHello(data []byte) (*tls.ClientHelloSpec, error) {
   fp := tls.Fingerprinter{AllowBluntMimicry: true}
   return fp.FingerprintClientHello(data)
}
