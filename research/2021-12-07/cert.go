package cert

import (
   "crypto/x509"
   "crypto/x509/pkix"
   "fmt"
   "math/big"
   "net"
   "time"
)

func newCert() {
   cert := &x509.Certificate{
   SerialNumber: big.NewInt(1658),
   Subject: pkix.Name{
   Organization:  []string{"Company, INC."},
   Country:       []string{"US"},
   Province:      []string{""},
   Locality:      []string{"San Francisco"},
   StreetAddress: []string{"Golden Gate Bridge"},
   PostalCode:    []string{"94016"},
   },
   IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
   NotBefore:    time.Now(),
   NotAfter:     time.Now().AddDate(10, 0, 0),
   SubjectKeyId: []byte{1, 2, 3, 4, 6},
   ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
   KeyUsage:     x509.KeyUsageDigitalSignature,
   }
   fmt.Println(cert)
}
