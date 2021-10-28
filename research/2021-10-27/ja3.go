package tls

import (
   "github.com/89z/parse/tls"
   "strconv"
   "strings"
   utls "github.com/refraction-networking/utls"
)

func parse(ja3 string) (*tls.ClientHello, error) {
   // This must be local, to prevent mutation.
   exts := make(map[string]utls.TLSExtension)
   tokens := strings.Split(ja3, ",")
   hello := tls.ClientHello{
      ClientHelloSpec: new(utls.ClientHelloSpec),
   }
   // Build CipherSuites.
   ciphers := strings.Split(tokens[1], "-")
   for _, cipher := range ciphers {
      cid, err := strconv.ParseUint(cipher, 10, 16)
      if err != nil {
         return nil, err
      }
      hello.CipherSuites = append(hello.CipherSuites, uint16(cid))
   }
   // Set extension 0.
   exts["0"] = &utls.SNIExtension{}
   // Set extension 10.
   var cids []utls.CurveID
   curves := strings.Split(tokens[3], "-")
   for _, curve := range curves {
      cid, err := strconv.ParseUint(curve, 10, 16)
      if err != nil {
         return nil, err
      }
      cids = append(cids, utls.CurveID(cid))
   }
   exts["10"] = &utls.SupportedCurvesExtension{cids}
   // Set extension 11.
   var pids []uint8
   points := strings.Split(tokens[4], "-")
   for _, point := range points {
      pid, err := strconv.ParseUint(point, 10, 8)
      if err != nil {
         return nil, err
      }
      pids = append(pids, uint8(pid))
   }
   exts["11"] = &utls.SupportedPointsExtension{pids}
   // Set extension 13. This cant be empty, so just use the Go default.
   exts["13"] = &utls.SignatureAlgorithmsExtension{
      []utls.SignatureScheme{
         0x804, 0x403, 0x807, 0x805, 0x806, 0x401,
         0x501, 0x601, 0x503, 0x603, 0x201, 0x203,
      },
   }
   // Set extension 16. If we leave this empty, then it will fail on any HTTP/2
   // servers.
   exts["16"] = &utls.ALPNExtension{
      []string{"http/1.1"},
   }
   // Set extension 23.
   exts["23"] = &utls.UtlsExtendedMasterSecretExtension{}
   // Set extension 35.
   exts["35"] = &utls.SessionTicketExtension{}
   // Set extension 65281.
   exts["65281"] = &utls.RenegotiationInfoExtension{}
   // Build extenions list.
   extensions := strings.Split(tokens[2], "-")
   for _, ext := range extensions {
      hello.Extensions = append(hello.Extensions, exts[ext])
   }
   // Return.
   ver, err := strconv.Atoi(tokens[0])
   if err != nil {
      return nil, err
   }
   hello.Version = uint16(ver)
   return &hello, nil
}
