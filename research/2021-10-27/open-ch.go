package main

import (
   "github.com/89z/parse/tls"
)

func value(ext tls.TLSExtension) (uint16, error) {
   data, err := io.ReadAll(ext)
   if err != nil {
      return 0, err
   }
   return binary.BigEndian.Uint16(data), nil
}

func marshalJA3(hello *tls.ClientHello) []byte {
   // An uint16 can contain numbers with up to 5 digits and an uint8 can
   // contain numbers with up to 3 digits, but we also need a byte for each
   // separating character, except at the end.
   var data []byte
   // Version
   data = strconv.AppendUint(data, uint64(hello.Version), 10)
   data = append(data, ',')
   // Cipher Suites
   if len(hello.CipherSuites) == 0 {
      data = append(data, ',')
   } else {
      for _, val := range hello.CipherSuites {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, '-')
      }
      // Replace last dash with a comma
      data[len(data)-1] = ','
   }
   // Extensions
   if len(hello.Extensions) == 0 {
      data = append(data, ',')
   } else {
      for _, ext := range hello.Extensions {
         val, err := value(ext)
         if err != nil {
            return nil
         }
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, '-')
      }
      // Replace last dash with a comma
      data[len(data)-1] = ','
   }
   // Elliptic curves
   if len(j.ellipticCurves) != 0 {
      for _, val := range j.ellipticCurves {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, '-')
      }
      // Replace last dash with a comma
      data[len(data)-1] = ','
   } else {
      data = append(data, ',')
   }
   // ECPF
   if len(j.ellipticCurvePF) != 0 {
      for _, val := range j.ellipticCurvePF {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, '-')
      }
      // Remove last dash
      data = data[:len(data)-1]
   }
   return data
}
