package ja3

import (
   "github.com/refraction-networking/utls"
   "strconv"
)

func supportedGroups(hello *tls.ClientHelloSpec) []tls.CurveID {
   for _, ext := range hello.Extensions {
      sc, ok := ext.(*tls.SupportedCurvesExtension)
      if ok {
         return sc.Curves
      }
   }
   return nil
}

func supportedPoints(hello *tls.ClientHelloSpec) []uint8 {
   for _, ext := range hello.Extensions {
      sp, ok := ext.(*tls.SupportedPointsExtension)
      if ok {
         return sp.SupportedPoints
      }
   }
   return nil
}

func bare(hello *tls.ClientHelloSpec) []byte {
   groups := supportedGroups(hello)
   points := supportedPoints(hello)
   // Version = uint16 => maximum = 65536 = 5chars + 1 field sep
   maxPossibleBufferLength := 5+1
   // CipherSuite = uint16 => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.CipherSuites)
   // uint16 = 2B => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(hello.Extensions)
   // uint16 = 2B => maximum = 65536 = 5chars
   maxPossibleBufferLength += (5+1)*len(groups)
   // uint8 = 1B => maximum = 256 = 3chars
   maxPossibleBufferLength += (3+1)*len(points)
   buffer := make([]byte, 0, maxPossibleBufferLength)
   buffer = strconv.AppendInt(buffer, int64(hello.TLSVersMin), 10)
   buffer = append(buffer, ',')
   // collect cipher suites
   last := len(hello.CipherSuites) - 1
   if len(hello.CipherSuites) > 1 {
      for _, cipher := range hello.CipherSuites[:last] {
         buffer = strconv.AppendInt(buffer, int64(cipher), 10)
         buffer = append(buffer, '-')
      }
   }
   // append last element if cipher suites are not empty
   if last != -1 {
      cipher := hello.CipherSuites[last]
      buffer = strconv.AppendInt(buffer, int64(cipher), 10)
   }
   buffer = append(buffer, ',')
   // collect extensions
   last = len(hello.Extensions) - 1
   if len(hello.Extensions) > 1 {
      for _, ext := range hello.Extensions[:last] {
         buffer = strconv.AppendInt(buffer, int64(ext), 10)
         buffer = append(buffer, '-')
      }
   }
   // append last element if extensions are not empty
   if last != -1 {
      ext := hello.Extensions[last]
      buffer = strconv.AppendInt(buffer, int64(ext), 10)
   }
   buffer = append(buffer, ',')
   // collect supported groups
   last = len(groups) - 1
   if len(groups) > 1 {
      for _, group := range groups[:last] {
         buffer = strconv.AppendInt(buffer, int64(group), 10)
         buffer = append(buffer, '-')
      }
   }
   // append last element if supported groups are not empty
   if last != -1 {
      buffer = strconv.AppendInt(buffer, int64(groups[last]), 10)
   }
   buffer = append(buffer, ',')
   // collect supported points
   last = len(points) - 1
   if len(points) > 1 {
      for _, point := range points[:last] {
         buffer = strconv.AppendInt(buffer, int64(point), 10)
         buffer = append(buffer, '-')
      }
   }
   // append last element if supported points are not empty
   if last != -1 {
      buffer = strconv.AppendInt(buffer, int64(points[last]), 10)
   }
   return buffer
}
