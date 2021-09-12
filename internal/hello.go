package main
import "github.com/refraction-networking/utls"

var cipherSuites = []uint16{
   0x1302, // TLS_AES_256_GCM_SHA384
   0x1303, // TLS_CHACHA20_POLY1305_SHA256
   0x1301, // TLS_AES_128_GCM_SHA256
   0xC02C, // TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
   0xC030, // TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
   0xC02B, // TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
   0xC02F, // TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
   0xCCA9, // TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
   0xCCA8, // TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
   0x009F, // TLS_DHE_RSA_WITH_AES_256_GCM_SHA384
   0x009E, // TLS_DHE_RSA_WITH_AES_128_GCM_SHA256
   0xCCAA, // TLS_DHE_RSA_WITH_CHACHA20_POLY1305_SHA256
   0xC0AF, // TLS_ECDHE_ECDSA_WITH_AES_256_CCM_8
   0xC0AD, // TLS_ECDHE_ECDSA_WITH_AES_256_CCM
   0xC0AE, // TLS_ECDHE_ECDSA_WITH_AES_128_CCM_8
   0xC0AC, // TLS_ECDHE_ECDSA_WITH_AES_128_CCM
   0xC024, // TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384
   0xC028, // TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384
   0xC023, // TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
   0xC027, // TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
   0xC00A, // TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
   0xC014, // TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
   0xC009, // TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
   0xC013, // TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
   0xC0A3, // TLS_DHE_RSA_WITH_AES_256_CCM_8
   0xC09F, // TLS_DHE_RSA_WITH_AES_256_CCM
   0xC0A2, // TLS_DHE_RSA_WITH_AES_128_CCM_8
   0xC09E, // TLS_DHE_RSA_WITH_AES_128_CCM
   0x006B, // TLS_DHE_RSA_WITH_AES_256_CBC_SHA256
   0x0067, // TLS_DHE_RSA_WITH_AES_128_CBC_SHA256
   0x0039, // TLS_DHE_RSA_WITH_AES_256_CBC_SHA
   0x0033, // TLS_DHE_RSA_WITH_AES_128_CBC_SHA
   0x009D, // TLS_RSA_WITH_AES_256_GCM_SHA384
   0x009C, // TLS_RSA_WITH_AES_128_GCM_SHA256
   0xC0A1, // TLS_RSA_WITH_AES_256_CCM_8
   0xC09D, // TLS_RSA_WITH_AES_256_CCM
   0xC0A0, // TLS_RSA_WITH_AES_128_CCM_8
   0xC09C, // TLS_RSA_WITH_AES_128_CCM
   0x003D, // TLS_RSA_WITH_AES_256_CBC_SHA256
   0x003C, // TLS_RSA_WITH_AES_128_CBC_SHA256
   0x0035, // TLS_RSA_WITH_AES_256_CBC_SHA
   0x002F, // TLS_RSA_WITH_AES_128_CBC_SHA
   0x00FF, // TLS_EMPTY_RENEGOTIATION_INFO_SCSV
}

// iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml
var preset = &tls.ClientHelloSpec{
   CipherSuites: cipherSuites,
   Extensions: []tls.TLSExtension{
      &tls.SNIExtension{},
      &tls.SupportedPointsExtension{
         SupportedPoints: []byte{0, 1, 2},
      },
      &tls.SupportedCurvesExtension{
         []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
            0x001E, // Curve448
            tls.CurveP521,
            tls.CurveP384,
         },
      },
      &tls.GenericExtension{Id: 35}, // session_ticket
      &tls.GenericExtension{Id: 22}, // encrypt_then_mac
      &tls.UtlsExtendedMasterSecretExtension{},
      &tls.SignatureAlgorithmsExtension{
         SupportedSignatureAlgorithms: []tls.SignatureScheme{
            tls.ECDSAWithP256AndSHA256,
            tls.PKCS1WithSHA256,
            tls.PSSWithSHA256,
         },
      },
      &tls.SupportedVersionsExtension{
         []uint16{tls.VersionTLS12},
      },
      &tls.PSKKeyExchangeModesExtension{
         []uint8{tls.PskModeDHE},
      },
      &tls.KeyShareExtension{
         []tls.KeyShare{
            {Group: tls.X25519},
         },
      },
   },
}
