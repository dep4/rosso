package main

import (
   "github.com/refraction-networking/utls"
)

func marshalJA3(hello *tls.ClientHelloSpec) []byte {
   // An uint16 can contain numbers with up to 5 digits and an uint8 can
   // contain numbers with up to 3 digits, but we also need a byte for each
   // separating character, except at the end.
   var data []byte
   // Version
   data = strconv.AppendUint(data, uint64(j.version), 10)
   data = append(data, commaByte)
   // Cipher Suites
   if len(j.cipherSuites) != 0 {
      for _, val := range j.cipherSuites {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, dashByte)
      }
      // Replace last dash with a comma
      data[len(data)-1] = commaByte
   } else {
      data = append(data, commaByte)
   }
   // Extensions
   if len(j.extensions) != 0 {
      for _, val := range j.extensions {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, dashByte)
      }
      // Replace last dash with a comma
      data[len(data)-1] = commaByte
   } else {
      data = append(data, commaByte)
   }
   // Elliptic curves
   if len(j.ellipticCurves) != 0 {
      for _, val := range j.ellipticCurves {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, dashByte)
      }
      // Replace last dash with a comma
      data[len(data)-1] = commaByte
   } else {
      data = append(data, commaByte)
   }
   // ECPF
   if len(j.ellipticCurvePF) != 0 {
      for _, val := range j.ellipticCurvePF {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, dashByte)
      }
      // Remove last dash
      data = data[:len(data)-1]
   }
   return data
}
