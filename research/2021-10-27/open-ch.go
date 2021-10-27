package tls

import (
   "encoding/binary"
   "github.com/89z/parse/tls"
   "io"
   "strconv"
   utls "github.com/refraction-networking/utls"
)

func value(ext utls.TLSExtension) (uint16, error) {
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
   var (
      supportedCurves []utls.CurveID
      supportedPoints []uint8
   )
   if len(hello.Extensions) == 0 {
      data = append(data, ',')
   } else {
      for _, iExt := range hello.Extensions {
         switch sExt := iExt.(type) {
         case *utls.SupportedCurvesExtension:
            supportedCurves = sExt.Curves
         case *utls.SupportedPointsExtension:
            supportedPoints = sExt.SupportedPoints
         }
         val, err := value(iExt)
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
   if len(supportedCurves) == 0 {
      data = append(data, ',')
   } else {
      for _, val := range supportedCurves {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, '-')
      }
      // Replace last dash with a comma
      data[len(data)-1] = ','
   }
   // ECPF
   if len(supportedPoints) > 0 {
      for _, val := range supportedPoints {
         data = strconv.AppendUint(data, uint64(val), 10)
         data = append(data, '-')
      }
      // Remove last dash
      data = data[:len(data)-1]
   }
   return data
}
